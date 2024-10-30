delimiter abcd
select "abcd" abcd
select 'abcd' abcd
delimiter "'efgh'"
select 'abcd' 'efgh'
select "'abcd'" 'efgh'
delimiter "`abcd`"
select '`abcd`' `abcd`
select "`abcd`" `abcd`
delimiter `"efgh"`
select '"abcd"' "efgh"
select "abcd" "efgh"
delimiter `"abcd";`
select '"abcd";' "abcd";
select '"aacd";' "abcd";
delimiter "`efgh`;"
select "`abcd`;" `efgh`;
select '`abcd`;' `efgh`;
delimiter "'abcd';"
select "'abcd';" 'abcd';
select "'abcd1';" 'abcd';
delimiter ab
select "ab'abcd';" ab
select "ab'abbd';" ab
select 'abcd'; ab
select 'ab'; ab