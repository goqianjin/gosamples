package bucketmodel

func NewCreateOption() *CreateOption {
	return &CreateOption{}
}

// ----

type CreateOption struct {
	Region string // z0, z1
}

func (opt *CreateOption) WithRegion(region string) *CreateOption {
	opt.Region = region
	return opt
}
