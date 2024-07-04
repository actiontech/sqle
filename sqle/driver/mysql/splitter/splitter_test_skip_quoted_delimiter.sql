delimiter abcd
select `abcd` abcd
select "abcd" abcd
select 'abcd' abcd
delimiter "'abcd'"
select `'abcd'` 'abcd'
select "'abcd'" 'abcd'
delimiter "`abcd`"
select '`abcd`' `abcd`
select "`abcd`" `abcd`
delimiter `"abcd"`
select '"abcd"' "abcd"
select `"abcd"` "abcd"
delimiter `"abcd";`
select '"abcd";' "abcd";
select `"abcd";` "abcd";
delimiter "`abcd`;"
select "`abcd`;" `abcd`;
select '`abcd`;' `abcd`;
delimiter "'abcd';"
select "'abcd';" 'abcd';
select `'abcd';` 'abcd';
delimiter ab
select "ab'abcd';" ab
select `ab'abcd';` ab
select `abcd`; ab