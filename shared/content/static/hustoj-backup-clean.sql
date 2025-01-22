delete from source_code
where
    solution_id in (
        select
            solution_id
        from
            solution
        where
            problem_id = 0
            and result > 4
    );

delete from source_code_user
where
    solution_id in (
        select
            solution_id
        from
            solution
        where
            problem_id = 0
            and result > 4
    );

delete from runtimeinfo
where
    solution_id in (
        select
            solution_id
        from
            solution
        where
            problem_id = 0
            and result > 4
    );

delete from compileinfo
where
    solution_id in (
        select
            solution_id
        from
            solution
        where
            problem_id = 0
            and result > 4
    );

update solution
set
    result = 5
where
    result < 4
    and in_date < curdate () - interval 3 day;

delete from solution
where
    problem_id = 0
    and result > 4;

--  cleanup trash from 6 month ago
delete from loginlog
where
    time < curdate () - interval 6 month;

delete from compileinfo
where
    solution_id < (
        select
            solution_id
        from
            solution
        where
            result = 11
            and in_date < curdate () - interval 6 month
        order by
            solution_id desc
        limit
            1
    );

delete from runtimeinfo
where
    solution_id < (
        select
            solution_id
        from
            solution
        where
            result = 11
            and in_date < curdate () - interval 6 month
        order by
            solution_id desc
        limit
            1
    );

repair table compileinfo,
contest,
contest_problem,
loginlog,
news,
privilege,
problem,
solution,
source_code,
users,
topic,
reply,
online,
sim,
mail;

optimize table compileinfo,
contest,
contest_problem,
loginlog,
news,
privilege,
problem,
solution,
source_code,
users,
topic,
reply,
online,
sim,
mail;
