package aws

import (
	"flag"
	"log"
	"os"
	"testing"
)

var FlagSweep bool
var SweeperFuncs map[string][]*Sweeper

type SweeperFunc func(i interface{}) error

func init() {
	flag.BoolVar(&FlagSweep, "sweep", false, "")
	SweeperFuncs = make(map[string][]*Sweeper)
}

type Sweeper struct {
	// Configuration for initializing the client connection for each Provider.
	// Ex google/config.go Config Struct
	// Ex aws/config.go Config Struct
	Config interface{}

	// Sweeper function that when invoked sweeps the Provider of specific
	// resources
	F SweeperFunc
}

func AddTestSweepers(name string, sf []*Sweeper) {
	if _, ok := SweeperFuncs[name]; ok {
		log.Printf("Error adding (%s) to SweeperFuncs: function already exists in map", name)
		os.Exit(1)
	}

	SweeperFuncs[name] = sf
}

func TestMain(m *testing.M) {
	flag.Parse()
	if FlagSweep {
		for n, s := range SweeperFuncs {
			log.Printf("[DEBUG] Running (%s) Sweeper...\n", n)
			for _, f := range s {
				client, err := f.Config.(*Config).Client()
				if err != nil {
					log.Printf("[ERR] Error with aws client: %s", err)
					os.Exit(1)
				}
				if err := f.F(client); err != nil {
					log.Printf("Error in (%s) Sweeper: %s", n, err)
					os.Exit(1)
				}
			}
		}
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestFake(t *testing.T) {
	log.Printf("Fake test")
	t.Fatalf("fall through")
}
