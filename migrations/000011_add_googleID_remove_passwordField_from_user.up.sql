alter table users add column google_id varchar(255) not null ;
alter table users drop column password;

create index indx_users_google_id on users(google_id); 