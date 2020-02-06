const user = process.env.MONGO_INITDB_ROOT_USERNAME
const pwd = process.env.MONGO_INITDB_ROOT_PASSWORD
const db = process.env.MONGO_INITDB_DATABASE

db.createUser({
    user,
    pwd,
    roles: [
        {
            role: "readWrite",
            db
        }
    ]
})