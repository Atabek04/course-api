-- +goose Up
CREATE TABLE IF NOT EXISTS department_info (
                             id BIGSERIAL PRIMARY KEY,
                            department_name varchar(255) not null,
                            staff_quantity int not null,
                            department_director varchar(255) not null,
                            module_id int not null ,
    CONSTRAINT fk_department
                                           FOREIGN KEY(module_id)
                                           REFERENCES module_info(id)

);

-- +goose Down
DROP TABLE IF EXISTS department_info;

