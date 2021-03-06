/*
Copyright DTCC 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

         http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package org.hyperledger.fabric.shim;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import com.google.protobuf.ByteString;

public class ChaincodeStub {
	
    private static Log logger = LogFactory.getLog(ChaincodeStub.class);
    
    private final String uuid;
    private final Handler handler;
    private final List<ByteString> args;

    public ChaincodeStub(String uuid, Handler handler, List<ByteString> args) {
        this.uuid = uuid;
        this.handler = handler;
        this.args = Collections.unmodifiableList(args);
    }

    public List<byte[]> getArgs() {
		return args
				.stream()
				.map(x -> x.toByteArray())
				.collect(Collectors.toList());
	}
    
    public List<String> getArgsAsStrings() {
		return args
				.stream()
				.map(x -> x.toStringUtf8())
				.collect(Collectors.toList());
	}
    
    /**
     * Gets the UUID of this stub
     *
     * @return the id used to identify this communication channel
     */
    public String getUuid() {
        return uuid;
    }

    /**
     * Get the state of the provided key from the ledger, and returns is as a string
     *
     * @param key the key of the desired state
     * @return the String value of the requested state
     */
    public String getState(String key) {
        return handler.handleGetState(key, uuid).toStringUtf8();
    }

    /**
     * Puts the given state into a ledger, automatically wrapping it in a ByteString
     *
     * @param key   reference key
     * @param value value to be put
     */
    public void putState(String key, String value) {
        handler.handlePutState(key, ByteString.copyFromUtf8(value), uuid);
    }

    /**
     * Deletes the state of the given key from the ledger
     *
     * @param key key of the state to be deleted
     */
    public void delState(String key) {
        handler.handleDeleteState(key, uuid);
    }
    /**
     * Given a start key and end key, this method returns a map of items with value converted to UTF-8 string.
     *
     * @param startKey
     * @param endKey
     * @return
     */
    //TODO: Uncomment and fix range query with new proto type
    /*
    public Map<String, String> getStateByRange(String startKey, String endKey) {
        Map<String, String> retMap = new HashMap<>();
        for (Map.Entry<String, ByteString> item : getStateByRangeRaw(startKey, endKey).entrySet()) {
            retMap.put(item.getKey(), item.getValue().toStringUtf8());
        }
        return retMap;
    }
    */
    /**
     * This method is same as getStateByRange, except it returns value in ByteString, useful in cases where
     * serialized object can be retrieved.
     *
     * @param startKey
     * @param endKey
     * @return
     */
    //TODO: Uncomment and fix range query with new proto type
    /*
    public Map<String, ByteString> getStateByRangeRaw(String startKey, String endKey) {
        Map<String, ByteString> map = new HashMap<>();
        for (ChaincodeShim.QueryStateKeyValue mapping : handler.handleGetStateByRange(
                startKey, endKey, uuid).getKeysAndValuesList()) {
            map.put(mapping.getKey(), mapping.getValue());
        }
        return map;
    }
    */

    /**
     * Given a partial composite key, this method returns a map of items (whose key's prefix 
     * matches the given partial composite key) with value converted to UTF-8 string and 
     * this methid should be used only for a partial composite key; For a full composite key, 
     * an iter with empty response would be returned.
	 *      
     * @param startKey
     * @param endKey
     * @return
     */

    //TODO: Uncomment and fix range query with new proto type
    /*
    public Map<String, String> getStateByPartialCompositeKey(String objectType, String[] attributes) {
        String partialCompositeKey = new String();
        partialCompositeKey = createCompositeKey(objectType, attributes);
        return getStateByRange(partialCompositeKey+"1", partialCompositeKey+":");
    }
    */

     /**
     * Given a set of attributes, this method combines these attributes to return a composite key. 
     *
     * @param objectType
     * @param attributes
     * @return
     */
    public String createCompositeKey(String objectType, String[] attributes) {
        String compositeKey = new String();
        compositeKey = compositeKey + objectType;
        for (String attribute : attributes) {
            compositeKey = compositeKey + attribute.length() + attribute;
        }
        return compositeKey;
    }

    /**
     * @param chaincodeName
     * @param function
     * @param args
     * @return
     */
    public String invokeChaincode(String chaincodeName, String function, List<ByteString> args) {
        return handler.handleInvokeChaincode(chaincodeName, function, args, uuid).toStringUtf8();
    }

    //------RAW CALLS------

    /**
     * @param key
     * @return
     */
    public ByteString getRawState(String key) {
        return handler.handleGetState(key, uuid);
    }

    /**
     * @param key
     * @param value
     */
    public void putRawState(String key, ByteString value) {
        handler.handlePutState(key, value, uuid);
    }

    /**
     *
     * @param startKey
     * @param endKey
     * @param limit
     * @return
     */
//	public GetStateByRangeResponse getStateByRangeRaw(String startKey, String endKey, int limit) {
//		return handler.handleGetStateByRange(startKey, endKey, limit, uuid);
//	}

    /**
     * Invokes the provided chaincode with the given function and arguments, and returns the
     * raw ByteString value that invocation generated.
     *
     * @param chaincodeName The name of the chaincode to invoke
     * @param function      the function parameter to pass to the chaincode
     * @param args          the arguments to be provided in the chaincode call
     * @return the value returned by the chaincode call
     */
    public ByteString invokeRawChaincode(String chaincodeName, String function, List<ByteString> args) {
        return handler.handleInvokeChaincode(chaincodeName, function, args, uuid);
    }
}
