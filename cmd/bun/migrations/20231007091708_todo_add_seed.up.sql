SET statement_timeout = 0;

--bun:split

insert into todo
(id, user_id, title, content, order_num, is_deleted, created_at, last_updated_at)
values (1, 0, '인생은 쓰다 하..', '오늘 집에 오다가 지하철에 지갑을 두고 내렸다', 1, false, '2022-10-10 11:30:30', '2022-10-10 11:30:30');
