delete from accounts where user_id in (select id from users where email like 'imi@raw.ninja%') ;
delete from users where email like 'imi@raw.ninja%';