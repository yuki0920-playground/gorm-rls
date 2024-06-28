package main

import (
	"context"
	"database/sql/driver"
	"fmt"
)

// connector, connがdriver.Connector, driver.SessionResetterを実装していることを確認
var (
	_ driver.Connector       = &connector{}
	_ driver.SessionResetter = &conn{}
)

type connector struct {
	dsn string
	d   driver.Driver
}

// Connect database/sqlにて、新しいconnectionが生成されるときに呼ばれる
func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	fmt.Println("Connect is called")

	cn, err := c.d.Open(c.dsn)
	if err != nil {
		return nil, fmt.Errorf("failed c.d.Open: %w", err)
	}
	if err = setApplicationName(ctx, cn); err != nil {
		return nil, fmt.Errorf("failed setApplicationName: %w", err)
	}
	return &conn{cn}, nil
}

func (c *connector) Driver() driver.Driver {
	return c.d
}

type conn struct {
	driver.Conn
}

// ResetSession database/sqlにて、ConnectionPoolからconnが取り出されるときに呼ばれる
func (c *conn) ResetSession(ctx context.Context) error {
	fmt.Println("ResetSession is called")

	if err := setApplicationName(ctx, c); err != nil {
		return fmt.Errorf("failed setApplicationName: %w", err)
	}
	return nil
}

func setApplicationName(ctx context.Context, cn driver.Conn) error {
	tenantID := ctx.Value("tenant_id")
	if tenantID == nil {
		fmt.Println("tenant_id is nil")
		return nil
	}

	fmt.Println("tenant_id:", tenantID)

	stmt, err := cn.Prepare(fmt.Sprintf("SET app.tenant_id = '%s'", tenantID))
	if err != nil {
		return fmt.Errorf("failed Prepare: %w", err)
	}

	args := []driver.NamedValue{}
	_, err = stmt.(driver.StmtExecContext).ExecContext(context.TODO(), args)
	if err != nil {
		return fmt.Errorf("failed stmt.Exec: %w", err)
	}

	return nil
}
