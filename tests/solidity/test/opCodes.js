// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

const TodoList = artifacts.require('./OpCodes.sol')
const assert = require('assert')
let contractInstance
const Web3 = require('web3');
const web3 = new Web3(new Web3.providers.HttpProvider('http://localhost:8545'));
// const web3 = new Web3(new Web3.providers.HttpProvider('http://localhost:9545'));

contract('OpCodes', (accounts) => {
   beforeEach(async () => {
      contractInstance = await TodoList.deployed()
   })
   it('Should run without errors the majorit of opcodes', async () => {
     await contractInstance.test()
     await contractInstance.test_stop()

   })

   it('Should throw invalid op code', async () => {
    try{
      await contractInstance.test_invalid()
    }
    catch(error) {
      console.error(error);
    }
   })

   it('Should revert', async () => {
    try{
      await contractInstance.test_revert()    }
    catch(error) {
      console.error(error);
    }
   })
})
