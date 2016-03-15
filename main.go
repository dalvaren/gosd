// To run:
// go get github.com/githubnemo/CompileDaemon
// CompileDaemon -command="./gosd"

package main

import "os"
import "fmt"
import "time"
import "strings"
import "strconv"
import "syscall"
import "os/signal"

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

type Settings struct {
  TryRefreshAmount      int
  TryFindServiceAmount  int
  TryFindServiceDelay   time.Duration
}

var ServiceSettings Settings
var ServiceMaps map[string]ServiceMap
var ServiceUpdater Updater
var LastCronTime time.Time

func main() {
  fmt.Println("Starting GOSD client")

  driver := DriverRedis{}
  // start
  currentName := Start("service-2", "http://localhost:8885", driver)
  // fmt.Println(currentName)

  // get
  Get()

  // finish
  Finish(currentName)

  // add services manually
  AddServiceManually("service-2", "http://localhost:8886")
  AddServiceManually("service-2", "http://localhost:8887")

  // remove service with URL
  DeleteServiceWithURL("http://localhost:8881")
  // DeleteServiceWithURL("http://localhost:8884")

  // IterateServiceRoute
  fmt.Println(IterateServiceRoute("service-2"))
  fmt.Println(IterateServiceRoute("service-2"))
  fmt.Println(IterateServiceRoute("service-2"))
  fmt.Println(IterateServiceRoute("service-2"))
  fmt.Println(IterateServiceRoute("service-2"))

  // waiting closing.
  WaitClosing()
  // for {
  //       fmt.Println("sleeping...")
  //       time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
  //   }
}

func AddServiceManually(name, url string) {
  serviceCacheEntry := ServiceCacheEntry{
    Name: name + "-" + time.Now().Format("20060102150405.99999999"),
    URL: url,
  }
  ServiceUpdater.ServiceCacheEntries = append(ServiceUpdater.ServiceCacheEntries, serviceCacheEntry)
  recalculateServiceMaps(ServiceUpdater)
}

func UpdateByCron() {
  if LastCronTime.Add(120 * time.Minute).Before(time.Now()) {
    Get()
  }
}

func WaitClosing() {
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  signal.Notify(c, syscall.SIGTERM)
  go func() {
      <-c
      fmt.Println("Finishing service: " + ServiceUpdater.Name)
      Finish(ServiceUpdater.Name)
      os.Exit(1)
  }()
}

func IterateServiceRoute(serviceBaseName string) string {
  attemptsNumber := ServiceSettings.TryFindServiceAmount
  for attemptsNumber > 0 {
    url := getNextURLForService(serviceBaseName)
    if url != "" {
      return url
    }
    attemptsNumber--
    time.Sleep(ServiceSettings.TryFindServiceDelay)
  }
  return ""
}

func Finish(currentName string) {
  // delete
  Delete(currentName)
}

func Delete(currentName string) {
  ServiceUpdater.Driver.Delete(currentName)
  Get()
}

func DeleteServiceWithURL(url string) {
  for _,serviceCacheEntry := range ServiceUpdater.ServiceCacheEntries {
    if serviceCacheEntry.URL == url {
      Delete(serviceCacheEntry.Name)
    }
  }
}

func Start(name, url string, driver Driver) string {
  // start settings
  ServiceSettings.TryRefreshAmount = 3
  ServiceSettings.TryFindServiceAmount = 5
  ServiceSettings.TryFindServiceDelay = 3 * time.Second
  if os.Getenv("TryRefreshAmount") != "" {
    param,_ := strconv.Atoi(os.Getenv("TryRefreshAmount"))
    ServiceSettings.TryRefreshAmount = param
  }
  if os.Getenv("TryFindServiceAmount") != "" {
    param,_ := strconv.Atoi(os.Getenv("TryFindServiceAmount"))
    ServiceSettings.TryRefreshAmount = param
  }
  if os.Getenv("TryFindServiceDelay") != "" {
    param,_ := strconv.Atoi(os.Getenv("TryFindServiceDelay"))
    ServiceSettings.TryFindServiceDelay = time.Duration(param) * time.Second
  }

  // start cron
  LastCronTime = time.Now()

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
  val := tryRefreshForNTimes(ServiceSettings.TryRefreshAmount)
  ServiceUpdater.ServiceCacheEntries = []ServiceCacheEntry{}
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
