package discovery

import (
	"os"
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ConfigSuite struct{}

func (s *ConfigSuite) TestIsLegalBackend(t sweet.T) {
	Expect(isLegalBackend("consul")).To(BeTrue())
	Expect(isLegalBackend("etcd")).To(BeTrue())
	Expect(isLegalBackend("zookeeper")).To(BeTrue())
	Expect(isLegalBackend("consulx")).To(BeFalse())
	Expect(isLegalBackend("zk")).To(BeFalse())
}

func (s *ConfigSuite) TestTimeConversion(t sweet.T) {
	c := &Config{
		RawDiscoveryTTL:      30,
		RawDiscoveryInterval: 15,
		DiscoveryHost:        "localhost",
		DiscoveryBackend:     "zookeeper",
	}

	Expect(c.PostLoad()).To(BeNil())
	Expect(c.DiscoveryTTL).To(Equal(time.Second * 30))
	Expect(c.DiscoveryInterval).To(Equal(time.Second * 15))
}

func (s *ConfigSuite) TestIllegalTTL(t sweet.T) {
	c := &Config{
		RawDiscoveryTTL:      15,
		RawDiscoveryInterval: 20,
		DiscoveryHost:        "localhost",
		DiscoveryBackend:     "zookeeper",
	}

	Expect(c.PostLoad()).To(Equal(ErrIllegalTTL))
}

func (s *ConfigSuite) TestDefaultHost(t sweet.T) {
	os.Clearenv()
	os.Setenv("HOST", "default-host")

	c := &Config{
		RawDiscoveryTTL:      30,
		RawDiscoveryInterval: 15,
		DiscoveryBackend:     "zookeeper",
	}

	Expect(c.PostLoad()).To(BeNil())
	Expect(c.DiscoveryHost).To(Equal("default-host"))
}

func (s *ConfigSuite) TestIllegalHost(t sweet.T) {
	os.Clearenv()

	c := &Config{
		RawDiscoveryTTL:      30,
		RawDiscoveryInterval: 15,
		DiscoveryBackend:     "zookeeper",
	}

	Expect(c.PostLoad()).To(Equal(ErrIllegalHost))
}
