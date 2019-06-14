# table-tail


For Tests not to fail you need to have a postgres db running:

    docker run --name tail-postgres -e POSTGRES_USER=root -e POSTGRES_DB=test -p 5432:5432 -d postgres