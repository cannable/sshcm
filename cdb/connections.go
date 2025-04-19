package cdb

func (conndb *ConnectionDB) Add(c *Connection) (int64, error) {
	// Do we have a nickname?
	if len(c.Nickname) < 1 {
		return -1, ErrAddNoNickname
	}

	// See if the nickname already exists
	if conndb.ExistsByProperty("nickname", c.Nickname) {
		return -1, ErrAddNicknameExists
	}

	// Try adding the connection
	result, err := conndb.connection.Exec(`
		INSERT INTO connections (
			nickname,
			host,
			user,
			description,
			args,
			identity,
			command,
			binary
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		)`,
		c.Nickname,
		c.Host,
		sqlNullableString(c.User),
		sqlNullableString(c.Description),
		sqlNullableString(c.Args),
		sqlNullableString(c.Identity),
		sqlNullableString(c.Command),
		sqlNullableString(c.Binary))

	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return id, err
}

func (conndb *ConnectionDB) Exists(id int64) bool {
	var check int

	err := conndb.connection.QueryRow("SELECT id FROM connections WHERE id = $1", id).Scan(&check)

	return err == nil
}

func (conndb *ConnectionDB) ExistsByProperty(property string, value string) bool {

	if !isValidProperty(property) {
		return false
	}

	var check int

	err := conndb.connection.QueryRow(`
		SELECT COUNT(nickname)
		FROM connections
		WHERE`+property+" = $1", value).Scan(&check)

	return err == nil
}

func (conndb *ConnectionDB) Get(id int64) (Connection, error) {
	var c Connection

	// Get connection details from DB
	row := conndb.connection.QueryRow(`
		SELECT
			id,
			nickname,
			host,
			user,
			description,
			args,
			identity,
			command,
			binary
		FROM connections
		WHERE id = $1
	`, id)

	c, err := marshallConnection(row, conndb.Defaults)

	return c, err
}

func (conndb *ConnectionDB) GetAll() ([]*Connection, error) {
	var cns []*Connection
	rows, err := conndb.connection.Query("SELECT id FROM connections")

	if err != nil {
		return cns, err
	}

	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			return cns, err
		}

		c, err := conndb.Get(id)

		if err != nil {
			return cns, err
		}

		cns = append(cns, &c)
	}

	return cns, err
}

func (conndb *ConnectionDB) GetByProperty(property string, value string) (Connection, error) {
	var c Connection

	if !isValidProperty(property) {
		return c, ErrPropertyInvalid
	}

	// Get connection details from DB
	row := conndb.connection.QueryRow(`
		SELECT
			id,
			nickname,
			host,
			user,
			description,
			args,
			identity,
			command,
			binary
		FROM connections
		WHERE `+property+" = $1", value)

	c, err := marshallConnection(row, conndb.Defaults)

	return c, err
}

func (conndb *ConnectionDB) Update(c *Connection) error {
	// Do we have an id?
	if c.Id < 1 {
		return ErrUpdateNoId
	}

	// Does the ID exist?
	if !conndb.Exists(c.Id) {
		return ErrUpdateIdNotExist
	}

	// Do we have a nickname?
	if len(c.Nickname) < 1 {
		return ErrUpdateNoNickname
	}

	// Try updating the connection
	_, err := conndb.connection.Exec(`
		UPDATE connections SET
			nickname = $2,
			host = $3,
			user = $4,
			description = $5,
			args = $6,
			identity = $7,
			command = $8,
			binary = $9
		WHERE id = $1
		`,
		c.Id,
		c.Nickname,
		c.Host,
		sqlNullableString(c.User),
		sqlNullableString(c.Description),
		sqlNullableString(c.Args),
		sqlNullableString(c.Identity),
		sqlNullableString(c.Command),
		sqlNullableString(c.Binary))

	return err
}
