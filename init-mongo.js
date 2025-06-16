db = db.getSiblingDB("exercise-3");
db.createUser({
  user: "bookuser",
  pwd: "bookpass",
  roles: [{ role: "readWrite", db: "exercise-3" }]
});
