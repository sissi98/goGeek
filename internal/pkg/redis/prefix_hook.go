package redis

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/frame/g"
)

type PrefixHook struct{}

var (
	once         = sync.Once{}
	prefix       string
	keysCommands = map[string]struct{}{
		"del":         {},
		"unlink":      {},
		"exists":      {},
		"touch":       {},
		"mget":        {},
		"sdiff":       {},
		"sinter":      {},
		"sunion":      {},
		"pfcount":     {},
		"rename":      {},
		"renamenx":    {},
		"sdiffstore":  {},
		"sinterstore": {},
		"sunionstore": {},
		"pfmerge":     {},
	}
	specialCommands = map[string]struct{}{
		"zdiff":                      {},
		"command":                    {},
		"client":                     {},
		"echo":                       {},
		"ping":                       {},
		"quit":                       {},
		"migrate":                    {},
		"randomkey":                  {},
		"bitop":                      {},
		"scan":                       {},
		"scantype":                   {},
		"blpop":                      {},
		"brpop":                      {},
		"brpoplpush":                 {},
		"rpoplpush":                  {},
		"lmove":                      {},
		"blmove":                     {},
		"smove":                      {},
		"xread":                      {},
		"xreadgroup":                 {},
		"xpendingext":                {},
		"xclaim":                     {},
		"xclaimjustid":               {},
		"xautoclaim":                 {},
		"xautoclaimjustid":           {},
		"bzpopmax":                   {},
		"bzpopmin":                   {},
		"zinter":                     {},
		"zinterstore":                {},
		"zrangestore":                {},
		"zunionstore":                {},
		"zunion":                     {},
		"zdiffstore":                 {},
		"bgrewriteaof":               {},
		"bgsave":                     {},
		"configget":                  {},
		"configresetstat":            {},
		"configset":                  {},
		"configrewrite":              {},
		"dbsize":                     {},
		"flushall":                   {},
		"flushallasync":              {},
		"flushdb":                    {},
		"flushdbasync":               {},
		"info":                       {},
		"lastsave":                   {},
		"save":                       {},
		"shutdown":                   {},
		"shutdownsave":               {},
		"shutdownnosave":             {},
		"slaveof":                    {},
		"time":                       {},
		"debugobject":                {},
		"readonly":                   {},
		"readwrite":                  {},
		"memoryusage":                {},
		"eval":                       {},
		"evalsha":                    {},
		"scriptexists":               {},
		"scriptflush":                {},
		"scriptkill":                 {},
		"scriptload":                 {},
		"publish":                    {},
		"pubsubchannels":             {},
		"pubsubnumsub":               {},
		"pubsubnumpat":               {},
		"clusterslots":               {},
		"clusternodes":               {},
		"clustermeet":                {},
		"clusterforget":              {},
		"clusterreplicate":           {},
		"clusterresetsoft":           {},
		"clusterresethard":           {},
		"clusterinfo":                {},
		"clusterkeyslot":             {},
		"clustergetkeysinslot":       {},
		"clustercountfailurereports": {},
		"clustercountkeysinslot":     {},
		"clusterdelslots":            {},
		"clusterdelslotsrange":       {},
		"clustersaveconfig":          {},
		"clusterslaves":              {},
		"clusterfailover":            {},
		"clusteraddslots":            {},
		"clusteraddslotsrange":       {},
		"auth":                       {},
		"authacl":                    {},
		"select":                     {},
		"swapdb":                     {},
		"clientsetname":              {},
		"mset":                       {},
		"msetnx":                     {},
	}
)

func SetPrefix(pre string) {
	prefix = pre
}

func (r *PrefixHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	//为了跑test不报错
	once.Do(func() {
		if prefix == "" {
			prefix = g.Cfg().GetString("app.name") + "-" + g.Cfg().GetString("app.env") + "_"
		}
	})
	args := cmd.Args()
	if len(args) > 1 {
		command := args[0].(string)
		if _, ok := keysCommands[command]; ok {
			for i := 1; i < len(args); i++ {
				args[i] = prefix + args[i].(string)
			}
		} else if _, ok := specialCommands[command]; ok {
			processSpecialCommandsArgs(args)
		} else {
			args[1] = prefix + args[1].(string)
		}
	}
	return ctx, nil
}

func (r *PrefixHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (r *PrefixHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}
func (r *PrefixHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func processSpecialCommandsArgs(args []interface{}) {
	command := args[0].(string)
	switch command {
	//以下参照自github.com/go-redis/redis/v8/commands
	//对于stream暂时未做处理
	case "zinter", "zunion", "zdiff", "zinterstore", "zunionstore", "zdiffstore":
		index := 2
		if command != "zinter" && command != "zunion" && command != "zdiff" {
			args[1] = prefix + args[1].(string)
			index = 3
		}
		for i := index; i < len(args); i++ {
			if v := args[i].(string); v == "withscores" || v == "weights" || v == "aggregate" {
				break
			}
			args[i] = prefix + args[i].(string)
		}
	case "migrate":
		args[3] = prefix + args[3].(string)
	case "bitop":
		for i := 2; i < len(args); i++ {
			args[i] = prefix + args[i].(string)
		}
	case "blpop", "brpop", "bzpopmax", "bzpopmin":
		for i := 1; i < len(args)-1; i++ {
			args[i] = prefix + args[i].(string)
		}
	case "rpoplpush", "brpoplpush", "lmove", "blmove", "smove", "zrangestore":
		args[1] = prefix + args[1].(string)
		args[2] = prefix + args[2].(string)
	case "mset", "msetnx":
		switch val := args[1].(type) {
		case string:
			for i := 1; i < len(args); i += 2 {
				args[i] = prefix + args[i].(string)
			}
		case []string:
			for i := 0; i < len(val); i += 2 {
				val[i] = prefix + val[i]
			}
		case map[string]interface{}:
			newVal := make(map[string]interface{}, len(val))
			for k, v := range val {
				newKey := prefix + k
				newVal[newKey] = v
			}
			args[1] = newVal
		}

	}
}
