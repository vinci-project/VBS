package main

import (
	"fmt"
	"os"
	"strings"
	"runtime"
	"strconv"
	"time"
	"log"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"

	"github.com/hashicorp/go-version"
	"github.com/go-redis/redis"

	// Speedtest
	"github.com/surol/speedtest-cli/speedtest"
)

type PartitionStat = disk.PartitionStat

type HardwareInfo struct {
	NodeId string
	OsVersion string
	GoVersion string
	RedisVersion string
	CpuType string
	CpuMhz float64
	MemorySize uint64
	FreeMemorySize uint64
	HarddiskModel string
	HarddiskSize uint64
}

func (node *storageNode) getHardwareInformation() (HardwareInfo,error){
	
	var availableSpace uint64;
	var redisVersion string;
	
	// Host Info
	hostInfo, _ := host.Info();

	hostId := hostInfo.HostID;

	osVersion := hostInfo.PlatformVersion;
	
	// Go Version
	goVersion := runtime.Version();

	// Redis Version
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, redisErr := client.Ping().Result();
	if redisErr != nil{
		redisVersion = "No Redis Found"
	} else {
		clientInfo,_ := client.Info().Result();

		clientInfoLine := strings.Split(clientInfo,"\r\n")

		for _, line := range clientInfoLine {
			if strings.Contains(line,"redis_version") {
				redisVersion = strings.Replace(line,"redis_version:","",-1)
				break
			}
		}
	}


	//CPU
	cpuInfo, _ := cpu.Info();
	cpuType := cpuInfo[0].ModelName;
	cpuMhz := cpuInfo[0].Mhz;

	// Memory
	memoryInfo, _ := mem.VirtualMemory()
	memorySize := memoryInfo.Total;
	freeMemorySize := memoryInfo.Available;

	// Storage
	storagePathDisk, err := disk.Usage(node.StoragePath)

	if err != nil {
		availableSpace = 0;
	} else {
		availableSpace = storagePathDisk.Free;
	}

	// Partitions
	partitions,_ := disk.Partitions(true);


	var selectedDisk PartitionStat;

	for _, disk := range partitions {
		if strings.Contains(node.StoragePath,strings.ToLower(disk.Mountpoint)) {
			selectedDisk = disk;
		}
	}

	// Harddisk
	// block, errS := ghw.Block();
	// if errS != nil {
	// 	fmt.Printf("Error getting block storage info: %v", errS)
	// }


	fmt.Println("Node ID:",hostId)
	fmt.Printf("-------------------------------------------------------\n") 
	fmt.Println("Os Version:",osVersion)
	fmt.Printf("-------------------------------------------------------\n") 
	fmt.Println("Go Version:",goVersion)
	fmt.Printf("-------------------------------------------------------\n") 
	fmt.Println("Redis Version:",redisVersion)
	fmt.Printf("-------------------------------------------------------\n") 
	fmt.Printf("CPU type:%v, Mhz:%v \n",cpuType,cpuMhz)
	fmt.Printf("-------------------------------------------------------\n") 
	fmt.Printf("Memory Size:%v, free memory size:%v \n",memorySize,freeMemorySize)
	fmt.Printf("-------------------------------------------------------\n") 
	fmt.Printf("Disk Model For Storage Path:%v \n",selectedDisk.Fstype)
	fmt.Printf("-------------------------------------------------------\n") 
	fmt.Printf("Disk Available For Storage Path:%v \n",availableSpace)
	fmt.Printf("-------------------------------------------------------\n") 

	node.NodeId = hostId;
	node.HardwareInfo = HardwareInfo{
		NodeId:hostId,
		OsVersion:osVersion,
		GoVersion:goVersion,
		RedisVersion:redisVersion,
		CpuType:cpuType,
		CpuMhz:cpuMhz,
		MemorySize:memorySize,
		FreeMemorySize:freeMemorySize,
		HarddiskModel:selectedDisk.Fstype,
		HarddiskSize:availableSpace,
	}

	return node.HardwareInfo,nil;

}

func (node *storageNode) checkRequirement(){

	fmt.Printf("Checking Hardware Requirements....\n") 
	fmt.Printf("-------------------------------------------------------\n") 

	// OS

	if runtime.GOOS == "windows" {
		osInfo := strings.Split(node.HardwareInfo.OsVersion," ")[0];
		fmt.Printf("Current OS %v\n", osInfo)
		fmt.Printf("Expected OS Windows 7+\n")
		currentOSVersion, _ := version.NewVersion(osInfo)
		expectedOSVersion, _ := version.NewVersion("7")
		fmt.Printf("OS Version Match %v\n", currentOSVersion.GreaterThan(expectedOSVersion))
		fmt.Printf("-------------------------------------------------------\n") 
		
	}

	// Golang
	currentGoVersion, _ := version.NewVersion(strings.Replace(node.HardwareInfo.GoVersion,"go","",-1))
	expectedGoVersion, _ := version.NewVersion("1.11")
	fmt.Printf("Current Go Version %v\n", strings.Replace(node.HardwareInfo.GoVersion,"go","",-1))
	fmt.Printf("Expected Go Version 1.11\n")
	fmt.Printf("Go Version Match %v\n", currentGoVersion.GreaterThan(expectedGoVersion))
	fmt.Printf("-------------------------------------------------------\n") 
	// Redis

	if node.HardwareInfo.RedisVersion == "No Redis Found" {
		fmt.Printf("No Redis Found\n")
	} else {
		currentRedisVersion, _ := version.NewVersion(node.HardwareInfo.RedisVersion);
		expectedRedisVersion, _ := version.NewVersion("4.0")
		fmt.Printf("Current Redis Version %v\n", node.HardwareInfo.RedisVersion)
		fmt.Printf("Expected Redis Version 4.0+\n")
		fmt.Printf("Redis Version Match %v\n", currentRedisVersion.GreaterThan(expectedRedisVersion))
	}
	fmt.Printf("-------------------------------------------------------\n") 

	// CPU type
	
	cpuInfo := strings.Split(node.HardwareInfo.CpuType," ");
	for _, cpuText := range cpuInfo {
		if strings.Contains(cpuText ,"GHz") {
			cpuSize,_:= strconv.ParseFloat(strings.Replace(cpuText,"GHz","",-1), 64)
			fmt.Printf("Current CPU Type %vGhz\n", cpuSize)
			fmt.Printf("Expected CPU Type 2Ghz+\n")
			fmt.Printf("CPU Type Match %v\n", cpuSize >= 2.0)
			break
		}
	}
	fmt.Printf("-------------------------------------------------------\n") 

	// Free Memory
	fmt.Printf("Current Free Memory %v\n", node.HardwareInfo.FreeMemorySize)
	fmt.Printf("Expected Free Memory 2GB+\n")
	fmt.Printf("Free Memory 2GB+ Match %v\n", node.HardwareInfo.FreeMemorySize > 2 * 1024 * 1024 * 1024)
	fmt.Printf("-------------------------------------------------------\n") 


	// Harddisk Available
	fmt.Printf("Current Harddisk Available %v\n", node.HardwareInfo.HarddiskSize)
	fmt.Printf("Expected Harddisk Available 1TB+\n")
	fmt.Printf("Harddisk Available Match %v\n", node.HardwareInfo.HarddiskSize > 1024 * 1024 * 1024 * 1024)
	fmt.Printf("-------------------------------------------------------\n") 


	// Storage Path Permission
	// info,_ := os.Stat(node.StoragePath);
	// mode:= info.Mode();
	// fmt.Println("Storage path mode is",info,mode)
	tempPath := node.StoragePath+"test.txt";
	f, writeErr := os.Create(tempPath);
	if writeErr !=nil {
		fmt.Printf("%v\n",writeErr) 
		fmt.Printf("Storage path has no permission\n") 
	} else {
		fmt.Printf("Storage path '%v' has permission\n",node.StoragePath) 
		f.Close();
		removeErr := os.Remove(tempPath);
		if removeErr != nil {
			fmt.Printf("%v\n",removeErr) 
		} else {
			fmt.Println("==> Done deleting Test file")
		}
	}
	
	fmt.Printf("-------------------------------------------------------\n") 

	// SpeedTest
	runTest()
}

func runTest(){
	fmt.Println("Doing Speedtest....")
	opts := new(speedtest.Opts)
	opts.Quiet = true;
	opts.Secure = true;
	client := speedtest.NewClient(opts)
	server := selectServer(client);
	downloadSpeed := server.DownloadSpeed()

	fmt.Printf("Current Bandwith %.2fMbits\n", float64(downloadSpeed) / (1 << 17))
	fmt.Printf("Expected Bandwith 10Mbits\n")
	fmt.Printf("Bandwith Available Match %v\n", (float64(downloadSpeed) / (1 << 18)) > 10.00 )
	fmt.Printf("-------------------------------------------------------\n") 

	// reportSpeed("Download", downloadSpeed)
	// uploadSpeed := server.UploadSpeed()
	// reportSpeed("Upload", uploadSpeed)
}

func reportSpeed(prefix string, speed int) {
	fmt.Printf("%s: %.2f Mib/s\n", prefix, float64(speed) / (1 << 17))
}

func selectServer(client speedtest.Client) (selected *speedtest.Server) {
	// if opts.Server != 0 {
	// 	servers, err := client.AllServers()
	// 	if err != nil {
	// 		log.Fatal("Failed to load server list: %v\n", err)
	// 		return nil
	// 	}
	// 	selected = servers.Find(opts.Server)
	// 	if selected == nil {
	// 		log.Fatalf("Server not found: %d\n", opts.Server)
	// 		return nil
	// 	}
	// 	selected.MeasureLatency(speedtest.DefaultLatencyMeasureTimes, speedtest.DefaultErrorLatency)
	// } else {
		servers, err := client.ClosestServers()
		if err != nil {
			log.Fatal("Failed to load server list\n")
			return nil
		}
		selected = servers.MeasureLatencies(
			speedtest.DefaultLatencyMeasureTimes,
			speedtest.DefaultErrorLatency).First()
	// }

	// if opts.Quiet {
		log.Printf("Ping: %d ms\n", selected.Latency / time.Millisecond)
	// } else {
		// client.Log("Hosted by %s (%s) [%.2f km]: %d ms\n",
		// 	selected.Sponsor,
		// 	selected.Name,
		// 	selected.Distance,
		// 	selected.Latency / time.Millisecond)
	// }

	return selected
}