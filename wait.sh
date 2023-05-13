echo "Waiting for 5 seconds..."
# sleep 50

function is_postgres_ready() {
    PGPASSWORD=${POSTGRES_PASSWORD} psql -h postgres -p 5432 -U ${POSTGRES_USER} -c "select 1;" > /dev/null 2>&1
    return $?
}
until is_postgres_ready
do
    echo "Waiting for PostgresSql to be ready..."
    sleep 1
done
echo "PostgresSql is ready."
migrate -path migrations -database ${DATABASE_URL} up
# go run .
./myapp

echo "Continuing with the Docker build process..."