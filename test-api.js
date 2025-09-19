// Test script to verify frontend API calls
async function testAPI() {
  try {
    console.log('Testing GetTasks...');
    const response = await fetch('http://localhost:8080/todo.v1.TodoService/GetTasks', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    const data = await response.json();
    console.log('GetTasks response:', data);
    
    console.log('\nTesting AddTask...');
    const addResponse = await fetch('http://localhost:8080/todo.v1.TodoService/AddTask', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ text: 'API测试任务' }),
    });
    
    const addData = await addResponse.json();
    console.log('AddTask response:', addData);
    
    console.log('\nTesting GetTasks again...');
    const finalResponse = await fetch('http://localhost:8080/todo.v1.TodoService/GetTasks', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    const finalData = await finalResponse.json();
    console.log('Final GetTasks response:', finalData);
    
  } catch (error) {
    console.error('API test failed:', error);
  }
}

testAPI();