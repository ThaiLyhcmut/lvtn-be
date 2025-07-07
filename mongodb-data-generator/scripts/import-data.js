const { MongoClient, ObjectId } = require('mongodb');
const fs = require('fs');
const path = require('path');
const { program } = require('commander');
const chalk = require('chalk');
const indexes = require('../config/indexes');

program
  .option('-u, --uri <uri>', 'MongoDB connection URI', indexes.mongoURI)
  .option('-d, --database <name>', 'Database name', 'lvtn')
  .option('-c, --collections <items>', 'Collections to import (comma separated)', '')
  .option('--drop', 'Drop existing collections before import')
  .option('--indexes', 'Create indexes after import', true)
  .parse(process.argv);

const options = program.opts();

console.log(chalk.blue('MongoDB Data Importer'));
console.log(chalk.gray('===================='));

const importData = async () => {
  const client = new MongoClient(options.uri);
  
  try {
    await client.connect();
    console.log(chalk.green('✓ Connected to MongoDB'));
    
    const db = client.db(options.database);
    const dataDir = path.join(__dirname, '..', 'data');
    
    // Get all JSON files or specific collections
    let files;
    if (options.collections) {
      files = options.collections.split(',').map(c => `${c.trim()}.json`);
    } else {
      files = fs.readdirSync(dataDir).filter(f => f.endsWith('.json'));
    }
    
    console.log(chalk.yellow(`\nImporting ${files.length} collections...`));
    
    for (const file of files) {
      const collectionName = path.basename(file, '.json');
      const filePath = path.join(dataDir, file);
      
      if (!fs.existsSync(filePath)) {
        console.log(chalk.red(`✗ File not found: ${filePath}`));
        continue;
      }
      
      const data = JSON.parse(fs.readFileSync(filePath, 'utf8'));
      
      if (data.length === 0) {
        console.log(chalk.gray(`  Skipping empty collection: ${collectionName}`));
        continue;
      }
      
      const collection = db.collection(collectionName);
      
      // Drop collection if requested
      if (options.drop) {
        try {
          await collection.drop();
          console.log(chalk.gray(`  Dropped existing collection: ${collectionName}`));
        } catch (err) {
          // Collection might not exist
        }
      }
      
      // Convert string dates to Date objects and ObjectId strings to ObjectId
      const processedData = data.map(doc => {
        const processDocument = (obj) => {
          const processed = { ...obj };
          Object.keys(processed).forEach(key => {
            // Convert _id and fields ending with _id to ObjectId
            if ((key === '_id' || key.endsWith('_id')) && processed[key]) {
              if (typeof processed[key] === 'string') {
                processed[key] = new ObjectId(processed[key]);
              } else if (processed[key].$oid) {
                processed[key] = new ObjectId(processed[key].$oid);
              }
            }
            // Convert ISO date strings to Date objects
            else if (typeof processed[key] === 'string' && 
                /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/.test(processed[key])) {
              processed[key] = new Date(processed[key]);
            }
            // Handle nested objects
            else if (processed[key] && typeof processed[key] === 'object' && !Array.isArray(processed[key])) {
              processed[key] = processDocument(processed[key]);
            }
            // Handle arrays of objects
            else if (Array.isArray(processed[key])) {
              processed[key] = processed[key].map(item => {
                if (typeof item === 'object' && item !== null) {
                  return processDocument(item);
                }
                return item;
              });
            }
          });
          return processed;
        };
        
        return processDocument(doc);
      });
      
      // Insert data
      const result = await collection.insertMany(processedData);
      console.log(chalk.green(`✓ Imported ${result.insertedCount} documents into ${collectionName}`));
      
      // Create indexes if requested
      if (options.indexes && indexes[collectionName]) {
        const indexSpecs = indexes[collectionName];
        for (const spec of indexSpecs) {
          await collection.createIndex(spec);
        }
        console.log(chalk.gray(`  Created ${indexSpecs.length} indexes for ${collectionName}`));
      }
    }
    
    // Display collection stats
    console.log(chalk.yellow('\nCollection Statistics:'));
    const collections = await db.listCollections().toArray();
    for (const col of collections) {
      const count = await db.collection(col.name).countDocuments();
      console.log(chalk.green(`  ${col.name}: ${count} documents`));
    }
    
  } catch (error) {
    console.error(chalk.red('✗ Error:'), error.message);
    process.exit(1);
  } finally {
    await client.close();
    console.log(chalk.blue('\n✓ Import completed'));
  }
};

importData();