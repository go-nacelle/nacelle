package nacelle

import (
	"github.com/aphistic/sweet"
	"github.com/efritz/nacelle/log"
	. "github.com/onsi/gomega"
)

type ServiceSuite struct{}

func (s *ServiceSuite) TestGetAndSet(t sweet.T) {
	container := NewServiceContainer()
	container.Set("a", &IntWrapper{10})
	container.Set("b", &FloatWrapper{3.14})
	container.Set("c", &IntWrapper{25})

	value1, err1 := container.Get("a")
	Expect(err1).To(BeNil())
	Expect(value1).To(Equal(&IntWrapper{10}))

	value2, err2 := container.Get("b")
	Expect(err2).To(BeNil())
	Expect(value2).To(Equal(&FloatWrapper{3.14}))

	value3, err3 := container.Get("c")
	Expect(err3).To(BeNil())
	Expect(value3).To(Equal(&IntWrapper{25}))
}

func (s *ServiceSuite) TestGetLogger(t sweet.T) {
	container := NewServiceContainer()
	logger, _ := log.NewGomolShim(&LoggingConfig{})
	err := container.Set("logger", logger)
	Expect(err).To(BeNil())
	Expect(container.GetLogger()).To(Equal(logger))
}

func (s *ServiceSuite) TestGetUnregisteredLogger(t sweet.T) {
	Expect(NewServiceContainer().GetLogger()).NotTo(BeNil())
}

func (s *ServiceSuite) TestSetBadLogger(t sweet.T) {
	Expect(NewServiceContainer().Set("logger", struct{}{})).To(Equal(ErrIllegalLogger))
}

func (s *ServiceSuite) TestInject(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	obj := &TestSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value.val).To(Equal(42))
}

func (s *ServiceSuite) TestInjectNonStruct(t sweet.T) {
	container := NewServiceContainer()
	obj := func() error { return nil }
	err := container.Inject(obj)
	Expect(err).To(BeNil())
}

func (s *ServiceSuite) TestInjectMissingService(t sweet.T) {
	container := NewServiceContainer()
	obj := &TestSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("no service registered to key"))
}

func (s *ServiceSuite) TestInjectBadType(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &FloatWrapper{3.14})
	obj := &TestSimpleProcess{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("field 'Value' cannot be assigned a value of type *nacelle.IntWrapper"))
}

func (s *ServiceSuite) TestInjectOptional(t sweet.T) {
	container := NewServiceContainer()
	obj := &TestOptionalServiceProcess{}
	err := container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value).To(BeNil())

	container.Set("value", &IntWrapper{42})
	err = container.Inject(obj)
	Expect(err).To(BeNil())
	Expect(obj.Value.val).To(Equal(42))
}

func (s *ServiceSuite) TestInjectBadOptional(t sweet.T) {
	container := NewServiceContainer()
	obj := &TestBadOptionalServiceProcess{}
	err := container.Inject(obj)
	Expect(err).To(MatchError("field 'Value' has an invalid optional tag"))
}

func (s *ServiceSuite) TestUnsettableFields(t sweet.T) {
	container := NewServiceContainer()
	container.Set("value", &IntWrapper{42})
	err := container.Inject(&TestUnsettableService{})
	Expect(err).To(MatchError("field 'value' can not be set"))
}

func (s *ServiceSuite) TestDuplicateRegistration(t sweet.T) {
	container := NewServiceContainer()
	err1 := container.Set("dup", struct{}{})
	err2 := container.Set("dup", struct{}{})
	Expect(err1).To(BeNil())
	Expect(err2).To(Equal(ErrDuplicateServiceKey))
}

func (s *ServiceSuite) TestGetUnregisteredKey(t sweet.T) {
	container := NewServiceContainer()
	_, err := container.Get("unregistered")
	Expect(err).To(Equal(ErrUnregisteredServiceKey))
}

func (s *ServiceSuite) TestMustSetPanics(t sweet.T) {
	Expect(func() {
		container := NewServiceContainer()
		container.MustSet("unregistered", struct{}{})
		container.MustSet("unregistered", struct{}{})
	}).To(Panic())
}

func (s *ServiceSuite) TestMustGetPanics(t sweet.T) {
	Expect(func() {
		NewServiceContainer().MustGet("unregistered")
	}).To(Panic())
}

//
// Processes

type (
	IntWrapper struct {
		val int
	}

	FloatWrapper struct {
		val float64
	}

	TestSimpleProcess struct {
		Value *IntWrapper `service:"value"`
	}

	TestUnsettableService struct {
		value *IntWrapper `service:"value"`
	}

	TestOptionalServiceProcess struct {
		Value *IntWrapper `service:"value" optional:"true"`
	}

	TestBadOptionalServiceProcess struct {
		Value *IntWrapper `service:"value" optional:"yup"`
	}
)
