
INSERT INTO tenants (id, name) VALUES ('tenant1', 'Tenant One');
INSERT INTO tenants (id, name) VALUES ('tenant2', 'Tenant Two');

INSERT INTO projects (id, name, tenant_id) VALUES ('project1', 'Project One', 'tenant1');
INSERT INTO projects (id, name, tenant_id) VALUES ('project2', 'Project Two', 'tenant1');

INSERT INTO projects (id, name, tenant_id) VALUES ('project3', 'Project Three', 'tenant2');
INSERT INTO projects (id, name, tenant_id) VALUES ('project4', 'Project Four', 'tenant2');
