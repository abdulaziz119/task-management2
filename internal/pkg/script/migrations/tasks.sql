CREATE TABLE IF NOT EXISTS tasks (
                       id SERIAL PRIMARY KEY,
                       project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                       name VARCHAR(255) NOT NULL,
                       description TEXT,
                       assigned_to INT REFERENCES users(id) ON DELETE SET NULL,
                       status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'in_progress', 'completed'
                       priority VARCHAR(50) DEFAULT 'medium', -- 'low', 'medium', 'high'
                       due_date DATE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP DEFAULT NULL
);

ALTER TABLE tasks OWNER TO postgres;