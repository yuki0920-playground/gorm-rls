ALTER TABLE projects ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_policy ON projects USING (tenant_id = current_setting('app.tenant_id'));
