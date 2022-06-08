package zentao

// Options is the struct for configurations of the zentao plugin.

// All fields are required for this version
type Options struct {
	Namespace             string                `validate:"required"`
	StorageClassName      string                `validate:"required"`
	PersistentVolume      PersistentVolume      `validate:"required"`
	PersistentVolumeClaim PersistentVolumeClaim `validate:"required"`
	Deployment            Deployment            `validate:"required"`
	Service               Service               `validate:"required"`
}

type PersistentVolume struct {
	ZentaoPVName     string `validate:"required"`
	ZentaoPVCapacity string `validate:"required"`

	MysqlPVName     string `validate:"required"`
	MysqlPVCapacity string `validate:"required"`
}

type PersistentVolumeClaim struct {
	ZentaoPVCName     string `validate:"required"`
	ZentaoPVCCapacity string `validate:"required"`

	MysqlPVCName     string `validate:"required"`
	MysqlPVCCapacity string `validate:"required"`
}

type Deployment struct {
	Name             string `validate:"required"`
	Replicas         int    `validate:"required"`
	Image            string `validate:"required"`
	MysqlPasswdName  string `validate:"required"`
	MysqlPasswdValue string `validate:"required"`
}

type Service struct {
	Name     string `validate:"required"`
	NodePort int    `validate:"required"`
}
