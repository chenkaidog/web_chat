package id_gen

import (
	"os"
	"strconv"
	"strings"
	"time"
	"web_chat/biz/util/ip"

	"github.com/bytedance/gopkg/lang/fastrand"
)

func init() {
	idgen = NewIDGenerator(100)
}

func NewLogID() string {
	return idgen.NewLogID()
}

var idgen *IDGenerator

type IDGenerator struct {
	traceIDPool <-chan string
	spanIDPool  <-chan string
	logIDPool   <-chan string
	stop        chan interface{}
}

func NewIDGenerator(maxSize int) *IDGenerator {
	stop := make(chan interface{})
	idgen := &IDGenerator{
		logIDPool: newLogIdPool(maxSize, stop),
		stop:      stop,
	}

	return idgen
}

func (idgen *IDGenerator) Stop() {
	select {
	case <-idgen.stop:
	default:
		close(idgen.stop)
	}
}

func (idgen *IDGenerator) NewLogID() string {
	return <-idgen.logIDPool
}

func newTraceIdPool(size int, stop chan interface{}) <-chan string {
	pool := make(chan string, size)

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				sb := strings.Builder{}
				sb.WriteString(strconv.FormatUint(uint64(time.Now().UnixMilli()), 36))
				sb.WriteString(ip.IPv4Hex())
				sb.WriteString(strconv.FormatUint(fastrand.Uint64(), 36))
				sb.WriteString(strconv.FormatInt(int64(os.Getpid()), 10))

				pool <- sb.String()
			}
		}
	}()

	return pool
}

func newLogIdPool(size int, stop chan interface{}) <-chan string {
	pool := make(chan string, size)

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				sb := strings.Builder{}
				sb.WriteString(strconv.FormatUint(uint64(time.Now().UnixMilli()), 36))
				sb.WriteString(ip.IPv4Hex())
				sb.WriteString(strconv.FormatUint(fastrand.Uint64(), 36))
				sb.WriteString(strconv.FormatUint(uint64(os.Getpid()), 10))

				pool <- sb.String()
			}
		}
	}()

	return pool
}
