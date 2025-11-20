-- migrate:up
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE,
    password VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    status BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
);

-- (username,password,kdsatker,role_id,status,keterangan,unit_kerja,unit_id,kd_instansiunitorg,pejabat_id)

-- migrate:down
DROP TABLE users;

