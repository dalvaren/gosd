// To run:
// go get github.com/githubnemo/CompileDaemon
// CompileDaemon -command="./gosd"

package main

import "fmt"
import "time"
import "strings"

type Driver interface {
  Start(name, url string) string
  Get() (map[string]string, error)
  Delete(currentName string)
}


type ServiceCacheEntry struct {
  Name  string
  URL   string
}

type Updater struct {
  Name          string
  State         string // expired or updated
  TTL           time.Time
  Driver        Driver
  ServiceCacheEntries  []ServiceCacheEntry
}

type ServiceMap struct {
  Index int
  URLs  []string
}

var ServiceMaps map[string]ServiceMap
var ServiceUpdater Updater

func main() {
  fmt.Println("Starting GOSD client")

  driver := DriverRedis{}
  // start
  currentName := Start("service-2", "http://localhost:8881", driver)
  // fmt.Println(currentName)

  // get
  Get()

  // finish
  Finish(currentName)

  // IterateServiceRoute
  fmt.Println(IterateServiceRoute("service-2"))
  fmt.Println(IterateServiceRoute("service-2"))
  fmt.Println(IterateServiceRoute("service-2"))
}

func IterateServiceRoute(serviceBaseName string) string {
  return getNextURLForService(serviceBaseName)
}

func Finish(currentName string) {
  // delete
  Delete(currentName)
}

func Delete(currentName string) {
  ServiceUpdater.Driver.Delete(currentName)
  // RedisClient.HDel("gosd", currentName)
}

func Start(name, url string, driver Driver) string {
  // start
  currentName := driver.Start(name, url)
  ServiceUpdater = Updater{
    Name: currentName,
    State: "expired",
    TTL:    time.Now(),
    Driver: driver,
  }

  return currentName
}

func Get() {
  val := tryRefreshForNTimes(3)
  for key,value := range val {
    // populate Updater
    serviceCacheEntry := ServiceCacheEntry{
      Name: key,
      URL: value,
    }
    ServiceUpdater.ServiceCacheEntries = append(ServiceUpdater.ServiceCacheEntries, serviceCacheEntry)
  }
  recalculateServiceMaps(ServiceUpdater)
}

func getNextURLForService(baseName string) string {
  serviceMaps, exists := ServiceMaps[baseName]
  if !exists {
    return ""
  }
  if serviceMaps.Index >= len(ServiceMaps[baseName].URLs) {
    serviceMaps.Index = 0
  }
  serviceMaps.Index++
  ServiceMaps[baseName] = serviceMaps
  return ServiceMaps[baseName].URLs[(serviceMaps.Index - 1)]
}

func recalculateServiceMaps(updater Updater) {
  if len(updater.ServiceCacheEntries) == 0 {
    return
  }
  serviceMaps := map[string]ServiceMap{}
  for _,serviceCacheEntry := range updater.ServiceCacheEntries {
    serviceMap := getServiceMap(getServiceBaseName(serviceCacheEntry.Name), serviceMaps)
    serviceMap.URLs = append(serviceMap.URLs, serviceCacheEntry.URL)
    serviceMaps[getServiceBaseName(serviceCacheEntry.Name)] = serviceMap
  }

  ServiceMaps = serviceMaps
}


func getServiceMap(baseName string, serviceMaps map[string]ServiceMap) ServiceMap {
  for key,serviceMap := range serviceMaps {
    if key == baseName {
      return serviceMap
    }
  }
  return ServiceMap{
    Index: 0,
    URLs:  []string{},
  }
}

func getServiceBaseName(name string) string {
  if index := strings.LastIndex(name, "-"); index > -1 {
    return name[0:index]
  }
  return name
}

func getNextServiceURL(name string) string {
  return ""
}

func tryRefreshForNTimes(n int) map[string]string {
  for n > 0 {
    val, err := ServiceUpdater.Driver.Get()
    if err != nil {
      n--
    } else {
      return val
    }
  }
  return map[string]string{}
}
