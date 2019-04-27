const redis = require('redis')
const client = redis.createClient(6379, 'redis') // port & host. See docker-compose.yml

client.on( 'connect', () => console.log('Redis client connected') )

client.on( 'error', (err) => console.log('Something went wrong ', err) )

console.log("About to set test value")
client.set('my-test-key', 'my test value', redis.print)

console.log("About to get test value")
client.get('my-test-key', (error, result) => {
  if (error) {
    console.log(error)
    throw error
  }
  console.log('GET result ->', result)
})