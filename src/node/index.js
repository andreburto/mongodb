const { MongoClient, ServerApiVersion } = require('mongodb');

// Get the MONGO_URL environment variable
const mongoUrl = process.env.MONGO_URI;

const dbName = "Counter";
const collectionName = "count";

console.log(`Connecting to ${mongoUrl}`);

// Create a MongoClient with a MongoClientOptions object to set the Stable API version
const client = new MongoClient(mongoUrl,  {
    serverApi: {
        version: ServerApiVersion.v1,
        strict: true,
        deprecationErrors: true,
        setTimeout: 10000,
    }
}
);

async function run() {
  try {
    const database = client.db(dbName);
    const countCollection = database.collection(collectionName);

    const clearResult = await countCollection.deleteMany({});
    
    // Create a document to insert
    const doc = {
      total: 0,
    }
    
    const result = await countCollection.insertOne(doc);

    // Print the ID of the inserted document
    console.log(`A document was inserted with the _id: ${result.insertedId}`);
  } finally {
     // Close the MongoDB client connection
    await client.close();
  }
}
// Run the function and handle any errors
run().catch(console.dir);
