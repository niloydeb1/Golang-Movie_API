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
	// FORGOT_PASSWORD refers to password forgot action
	FORGOT_PASSWORD = USER_UPDATE_ACTION("forgot_password")
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
	// CREATE_USER refers to create user by admin
	CREATE_USER = USER_REGISTRATION_ACTION("create_user")
)

// ROLE role string
type ROLE string

const (
	// ADMIN refers to admin role
	ADMIN = ROLE("ADMIN")
	// VIEWER refers to user role
	VIEWER = ROLE("VIEWER")
)

// RESOURCE resource string
type RESOURCE string

const (
	// USER refers to user resource
	USER = RESOURCE("user")
	// MOVIE refers to movie resource
	MOVIE = RESOURCE("movie")
)

// PERMISSION permission string
type PERMISSION string

const (
	// CREATE refers to CREATE permission
	CREATE = PERMISSION("CREATE")
	// READ refers to READ permission
	READ = PERMISSION("READ")
	// UPDATE refers to UPDATE permission
	UPDATE = PERMISSION("UPDATE")
	// DELETE refers to DELETE permission
	DELETE = PERMISSION("DELETE")
)