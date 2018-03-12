package blah
//package main
//
//import (
//	"fmt"
//	"github.com/jessevdk/go-flags"
//	"github.com/op/go-logging"
//	"github.com/textileio/textile-go/cmd"
//	"github.com/textileio/textile-go/core"
//	"github.com/textileio/textile-go/mobile"
//	lockfile "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo/lock"
//	"os"
//	"os/signal"
//	"path/filepath"
//)
//
//var log = logging.MustGetLogger("main")
//
//type Opts struct {
//	Version bool `short:"v" long:"version" description:"Print the version number and exit"`
//}
//type Stop struct{}
//type Restart struct{}
//
//var stopServer Stop
//var restartServer Restart
//
//var opts Opts
//
//var parser = flags.NewParser(&opts, flags.Default)
//
//func main() {
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, os.Interrupt)
//	go func() {
//		for sig := range c {
//			log.Noticef("Received %s\n", sig)
//			//fmt.Println("Received interrupt signal, shutting down...")
//			//fmt.Println("(Hit ctrl-c again to force-shutdown the daemon.)")
//			if core.Node != nil {
//				//if core.Node.MessageRetriever != nil {
//				//	close(core.Node.MessageRetriever.DoneChan)
//				//	core.Node.MessageRetriever.Wait()
//				//}
//				//core.OfflineMessageWaitGroup.Wait()
//				core.PublishLock.Lock()
//				core.Node.Datastore.Close()
//				repoLockFile := filepath.Join(core.Node.RepoPath, lockfile.LockFile)
//				os.Remove(repoLockFile)
//				//core.Node.Wallet.Close()
//				core.Node.IpfsNode.Close()
//			}
//			os.Exit(1)
//		}
//	}()
//	//parser.AddCommand("init",
//	//	"initialize a new repo and exit",
//	//	"Initializes a new repo without starting the server",
//	//	&cmd.Init{})
//	//parser.AddCommand("status",
//	//	"get the repo status",
//	//	"Returns the status of the repo ― Uninitialized, Encrypted, Decrypted. Also returns whether Tor is available.",
//	//	&cmd.Status{})
//	//parser.AddCommand("setapicreds",
//	//	"set API credentials",
//	//	"The API password field in the config file takes a SHA256 hash of the password. This command will generate the hash for you and save it to the config file.",
//	//	&cmd.SetAPICreds{})
//	parser.AddCommand("start",
//		"start the textile server",
//		"The start command starts the textile server",
//		&cmd.Start{})
//	parser.AddCommand("stop",
//		"shutdown the server and disconnect",
//		"The stop command disconnects from peers and shuts down textile server",
//		&stopServer)
//	parser.AddCommand("restart",
//		"restart the server",
//		"The restart command shuts down the server and restarts",
//		&restartServer)
//	//parser.AddCommand("encryptdatabase",
//	//	"encrypt your database",
//	//	"This command encrypts the database containing your bitcoin private keys, identity key, and contracts",
//	//	&cmd.EncryptDatabase{})
//	//parser.AddCommand("decryptdatabase",
//	//	"decrypt your database",
//	//	"This command decrypts the database containing your bitcoin private keys, identity key, and contracts.\n [Warning] doing so may put your bitcoins at risk.",
//	//	&cmd.DecryptDatabase{})
//	//parser.AddCommand("restore",
//	//	"restore user data",
//	//	"This command will attempt to restore user data (profile, listings, ratings, etc) by downloading them from the network. This will only work if the IPNS mapping is still available in the DHT. Optionally it will take a mnemonic seed to restore from.",
//	//	&cmd.Restore{})
//	//parser.AddCommand("convert",
//	//	"convert this node to a different coin type",
//	//	"This command will convert the node to use a different cryptocurrency",
//	//	&cmd.Convert{})
//	parser.AddCommand("mobile",
//		"start the textile mobile server",
//		"The start command starts the textile mobile server",
//		&mobile.Start{})
//	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
//		fmt.Println(core.VERSION)
//		return
//	}
//	if _, err := parser.Parse(); err != nil {
//		os.Exit(1)
//	}
//}