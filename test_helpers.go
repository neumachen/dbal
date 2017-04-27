package dbal

var badSQLStmnt = `SSSSSS`
var insertCustomer = `insert into customers(first_name, last_name, address) values(:first_name, :last_name, CAST(NULLIF(:address, '') as jsonb));`
var selectCustomer = `select first_name, last_name, address from customers where first_name = :first_name AND last_name = :last_name;`
var selectAllCustomers = `select first_name, last_name, address from customers;`
var deleteCustomer = `delete from customers where first_name = :first_name and last_name = :last_name;`
