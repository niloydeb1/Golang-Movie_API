package enums

// ENVIRONMENT run environment
type ENVIRONMENT string

const (
	// PRODUCTION production environment
	PRODUCTION = ENVIRONMENT("PRODUCTION")
	// DEVELOP development environment
	DEVELOP = ENVIRONMENT("DEVELOP")
	// TEST test environment
	TEST = ENVIRONMENT("TEST")
)

const (
	// MONGO mongo as db
	MONGO = "MONGO"
	// INMEMORY in memory storage as db
	INMEMORY = "INMEMORY"
)

// USER_UPDATE_ACTION users update action
type USER_UPDATE_ACTION string

const (
	// RESET_PASSWORD refers to password reset action
	RESET_PASSWORD = USER_UPDATE_ACTION("reset_password")
	// UPDATE_STATUS refers to status update action
	UPDATE_STATUS = USER_UPDATE_ACTION("update_status")
)

// STATUS status update action
type STATUS string

const (
	// ACTIVE user status for active user
	ACTIVE = STATUS("active")
	// INACTIVE user status for inactive user
	INACTIVE = STATUS("inactive")
	// DELETED user status for deleted user
	DELETED = STATUS("deleted")
)

// USER_REGISTRATION_ACTION user registration action
type USER_REGISTRATION_ACTION string

const (
	// CREATE_ADMIN refers to create admin by superadmin
	CREATE_ADMIN = USER_REGISTRATION_ACTION("create_admin")
)

// ROLE role string
type ROLE string

const (
	// SUPERADMIN refers to superadmin role
	SUPERADMIN = ROLE("SUPERADMIN")
	// ADMIN refers to admin role
	ADMIN = ROLE("ADMIN")
	// USER refers to user role
	USER = ROLE("USER")
)