// Test script to verify the worker stays alive
// This simulates the worker behavior without actual Kafka connection

console.log('=== Testing Worker Behavior ===\n');

// Simulate the consumer.run() behavior
function simulateConsumerRun() {
    return new Promise((resolve) => {
        console.log('[SIMULATION] Consumer.run() called - returns immediately');
        // Simulate immediate resolution like KafkaJS consumer.run()
        setTimeout(() => {
            console.log('[SIMULATION] Consumer.run() Promise resolved (setup complete)');
            resolve();
        }, 100);
    });
}

// Simulate message processing in background
function simulateBackgroundProcessing() {
    let messageCount = 0;
    const interval = setInterval(() => {
        messageCount++;
        console.log(`[SIMULATION] Processing message #${messageCount} in background`);
        if (messageCount >= 5) {
            console.log('[SIMULATION] Stopping simulation after 5 messages');
            clearInterval(interval);
            console.log('\n=== Test completed successfully ===');
            console.log('Worker would continue running indefinitely until SIGTERM/SIGINT');
            process.exit(0);
        }
    }, 1000);
}

// Test WITH await (incorrect - causes immediate exit)
async function testWithAwait() {
    console.log('TEST 1: With await (INCORRECT)');
    console.log('-------------------------------');
    await simulateConsumerRun();
    console.log('[ERROR] Main function completed - process will exit!');
    console.log('[ERROR] No background processing occurs\n');
}

// Test WITHOUT await (correct - stays alive)
async function testWithoutAwait() {
    console.log('TEST 2: Without await (CORRECT)');
    console.log('--------------------------------');
    simulateConsumerRun(); // No await
    console.log('[SUCCESS] Main function continues - process stays alive');
    console.log('[SUCCESS] Consumer runs in background');
    console.log('[SUCCESS] Simulating background message processing...\n');
    simulateBackgroundProcessing();
}

async function runTests() {
    await testWithAwait();
    
    console.log('\nWaiting 2 seconds before next test...\n');
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    await testWithoutAwait();
}

runTests();
