package main

import (
	"flag"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"runtime/debug"
	"syscall"
	"tc/bootstrap"
	"time"
)

var config struct {
	net           string
	laddr         string
	queueSize     int
	dumpBatchSize int
	dumpInterval  time.Duration
	db            string
	tableCap      int
	enableProfile bool
	verbosity     int

	mongoStats           string
	mongoStatsCollection string

	mongoMain string

	debug struct {
		disableBalance bool
		disableLimits  bool
	}
}

var db *sqlx.DB
var mongo struct {
	stats *mgo.Database
	main  *mgo.Database
}

const (
	VerbosityNone = iota
	VerbosityAll
)

func init() {
	parseFlags()
	connectToDB()
	connectToMongo()
}

func parseFlags() {
	flag.StringVar(&config.net, "net", "tcp", "net type: tcp|udp|unix")
	flag.StringVar(&config.laddr, "l", ":8880", "address to listen")
	flag.IntVar(&config.queueSize, "queue-size", 1000, "input queue size")
	flag.IntVar(&config.dumpBatchSize, "dump-batch-size", 1000, "num of stats per aggregation")
	flag.DurationVar(&config.dumpInterval, "dump-interval", time.Minute, "run aggregation at interval")
	flag.StringVar(&config.db, "db", "postgres://tubecontext:PQ9kzeRU@d1444.webazilla.com:9191/tubecontext?sslmode=disable", "db connection string")
	flag.IntVar(&config.tableCap, "table-cap", 1000, "default table cap")
	flag.BoolVar(&config.enableProfile, "profile", false, "enable http profile")

	flag.StringVar(&config.mongoStats, "mongo-stats", "d1444.webazilla.com/rawstats", "mongo stats connection string")
	flag.StringVar(&config.mongoStatsCollection, "mongo-stats-collection", "rawstats", "mongo stats collection")

	flag.StringVar(&config.mongoMain, "mongo-main", "d2883.webazilla.com:27017/tubecontext", "mongo main connection string")

	flag.IntVar(&config.verbosity, "v", config.verbosity, "show more information during execution")
	flag.BoolVar(&config.debug.disableBalance, "disable-balance", false, "disable balance calc")
	flag.BoolVar(&config.debug.disableLimits, "disable-limits", false, "disable limits updates")
	flag.Parse()
}

func connectToDB() {
	db = sqlx.MustConnect("postgres", config.db)
}

func connectToMongo() {
	connect := func(url string) *mgo.Database {
		ss, err := mgo.Dial(url)
		checkErr(err, "mongo connection failed")
		return ss.DB("")
	}
	mongo.main = connect(config.mongoMain)
	mongo.stats = connect(config.mongoStats)

	mongo.stats.C(config.mongoStatsCollection).EnsureIndexKey("-date")
}

func main() {
	// log.SetFlags(log.Flags() | log.Lshortfile)
	if config.enableProfile {
		go runProfile()
	}
	bootstrap.Run(func() {
		server := newServer()
		signal.Notify(server.runDump, syscall.SIGUSR1)
		server.run()
	})
}

func runProfile() {
	log.Println(http.ListenAndServe("localhost:6062", nil))
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, "::", err)
	}
}

func checkMongoErr(err error, mongo *mgo.Database) error {
	if err != nil {
		mongo.Session.Refresh()
		log.Println(err)
	}
	return err
}

func catchPanic() {
	if r := recover(); r != nil {
		debug.PrintStack()
		log.Println("PANIC:", r)
		return
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
