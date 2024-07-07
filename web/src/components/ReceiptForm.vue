<template>
  <div className="form">
    <h1>Totals and Subtotals</h1>
    <form id="receipt-form" class="receipt-form">
      <div class="totals">
        <div class="receipt-field">
          <label class="receipt-field-label">Total: </label>
          <input type="text" v-model.number="total" />
        </div>
        <div class="receipt-field">
          <label class="receipt-field-label">Subtotal: </label>
          <input type="text" v-model.number="subtotal" />
        </div>
      </div>
      <h1>Line Items</h1>
      <div v-for="(_, index) of Array(lineItemCount)" class="line-item">
        <div class="receipt-field">
          <label class="receipt-field-label">Item: </label>
          <input type="text" v-model="lineItems[index].item" />
        </div>
        <div class="receipt-field">
          <label class="receipt-field-label">Person: </label>
          <input type="text" v-model="lineItems[index].person" />
        </div>
        <div class="receipt-field">
          <label class="receipt-field-label">Cost: </label>
          <input type="text" v-model.number="lineItems[index].cost" />
        </div>
      </div>
      <button @click="addEmptyLineItem">Add Line Item</button>
      <button @click="generateSplits">Submit</button>
    </form>
    <div v-show="showResults" class="results-section">
      <h1 class="results-header">Results:</h1>
      <div v-for="split in splits" class="split">
        <h1 class="split-person">
          <u>{{ split.person }}</u> owes <u>${{ split.cost.toFixed(2) }}</u>
        </h1>
      </div>
    </div>
  </div>
</template>

<style scoped>
.receipt-field {
  margin-top: 2%;
  margin-bottom: 2%;
  width: 100%;
}

.receipt-field-label {
  font-size: 1.2em;
}

.line-item {
  border: solid 1px;
  margin-top: 1%;
  margin-bottom: 4%;
  padding-left: 2%;
  padding-right: 2%;
  padding-bottom: 2%;
  border-radius: 5%;
}

.split {
}

.split-person {
}

.totals {
  border: solid 1px;
  margin-top: 1%;
  margin-bottom: 4%;
  padding-left: 2%;
  padding-right: 2%;
  padding-bottom: 2%;
  border-radius: 5%;
}

.results-header {
  padding-top: 5%;
}

.results-section {
  padding-bottom: 6rem;
}

u {
  color: green;
}

h1 {
  font-size: 2em;
}

input {
  width: 100%;
  font-size: 1em;
  border-radius: 20px;
  border: solid 1px grey;
  padding-left: 5%;
  height: 30px;
}

button {
  width: 100%;
  font-size: 1.25em;
  margin-top: 5%;
  margin-bottom: 5%;
}
</style>

<script setup>
import { ref, toRaw } from "vue";

const lineItemCount = ref(1);
const total = ref(0);
const subtotal = ref(0);
const lineItems = ref([{}]);
const splits = ref([]);
const showResults = ref(false);

function addEmptyLineItem(event) {
  event.preventDefault();
  console.log("hello world");
  lineItemCount.value++;
  console.log(lineItemCount.value);
  lineItems.value.push({});
}

function updateLineItem(event, i, val, field) {
  event.preventDefault();
  lineItems[i][field] = val;
}

async function fetchSplits(request) {
  const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:3333";
  const response = await fetch(`${apiUrl}/split`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    console.log(response);
    throw new Error("Network response was not ok");
  }

  return response.json();
}

async function generateSplits(event) {
  event.preventDefault();
  const request = {
    total_cost: total.value,
    subtotal: subtotal.value,
    line_items: lineItems.value,
  };

  splits.value = [];

  try {
    const response = await fetchSplits(request);
    splits.value = response.people;
    showResults.value = true;
  } catch (error) {
    console.error("There was an error with the fetch operation:", error);
  }
}
</script>
