package guardian

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/nxadm/tail"
	"github.com/patrickmn/go-cache"
)

type Guardian struct {
	cache   *cache.Cache
	logPath string
}

func (g *Guardian) SaveList(path string) {
	g.cache.SaveFile(path)
}

func NewGuardian(cache *cache.Cache, logPath string) *Guardian {
	return &Guardian{
		cache:   cache,
		logPath: logPath,
	}
}

func (g *Guardian) StartProc(order string, args []string) bool {
	logFile := g.logPath + order + ".log"
	command := order
	for _, v := range args {
		command += " " + v
	}
	if pids, ok := g.getPid(command); ok {
		if len(pids) != 0 {
			fmt.Println("cmd already running , kill first pid:", pids[0])
			return false
		}
	}
	cmd := exec.Command("bash", "-c", "nohup "+command+" > "+logFile+" 2>&1 &")
	if err := cmd.Start(); err != nil {
		return false
	}
	if pids, ok := g.getPid(command); ok {
		if pids != nil {
			fmt.Println("CMD:        " + order)
			fmt.Println("PID:        " + pids[len(pids)-1])
			fmt.Println("STATUS:     " + "Start OK ")
			g.cache.SetDefault(order, map[string]string{
				"pid": pids[len(pids)-1],
				"cmd": command,
				"log": logFile,
			})
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (g *Guardian) ShowList() {
	for k, v := range g.cache.Items() {
		item := v.Object.(map[string]string)
		fmt.Printf("cmd: %s     pid: %s     args: %s \n", k, item["pid"], item["cmd"])
	}
}

func (g *Guardian) StopProc(order string) {
	if item, ok := g.cache.Get(order); ok {
		pid := item.(map[string]string)["pid"]
		if err := exec.Command("bash", "-c", "kill -9 "+pid).Start(); err != nil {
			fmt.Println("kill fail pid:", pid)
		} else {
			g.cache.Delete(order)
			fmt.Println("kill pid:" + pid + " ok")
		}
	} else {
		fmt.Println("kill fail pid not found")
	}
}

func (g *Guardian) ShowLog(order string) {
	if item, ok := g.cache.Get(order); ok {
		proc := item.(map[string]string)
		if t, err := tail.TailFile(proc["log"], tail.Config{Follow: true}); err != nil {
			fmt.Println("open file fail")
			return
		} else {
			for line := range t.Lines {
				fmt.Println(line.Text)
			}
		}
	} else {
		fmt.Println("cmd not found")
	}
}

func (g *Guardian) getPid(command string) ([]string, bool) {
	cmd := exec.Command("bash", "-c", "ps -aux | grep '"+command+"' | grep -v grep | grep -v gdn | awk '{print $2}'")
	if res, err := cmd.CombinedOutput(); err != nil {
		return nil, false
	} else {
		pids := strings.Split(string(res), "\n")
		return pids[:len(pids)-1], true
	}
}
