/**
 * ORY Hydra
 * Welcome to the ORY Hydra HTTP API documentation. You will find documentation for all HTTP APIs here.
 *
 * OpenAPI spec version: latest
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 2.4.5
 *
 * Do not edit the class manually.
 *
 */

(function(root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as an anonymous module.
    define(['ApiClient', 'model/JSONWebKey'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'), require('./JSONWebKey'));
  } else {
    // Browser globals (root is window)
    if (!root.OryHydra) {
      root.OryHydra = {};
    }
    root.OryHydra.SwaggerJwkUpdateSetKey = factory(root.OryHydra.ApiClient, root.OryHydra.JSONWebKey);
  }
}(this, function(ApiClient, JSONWebKey) {
  'use strict';




  /**
   * The SwaggerJwkUpdateSetKey model module.
   * @module model/SwaggerJwkUpdateSetKey
   * @version latest
   */

  /**
   * Constructs a new <code>SwaggerJwkUpdateSetKey</code>.
   * @alias module:model/SwaggerJwkUpdateSetKey
   * @class
   * @param kid {String} The kid of the desired key in: path
   * @param set {String} The set in: path
   */
  var exports = function(kid, set) {
    var _this = this;


    _this['kid'] = kid;
    _this['set'] = set;
  };

  /**
   * Constructs a <code>SwaggerJwkUpdateSetKey</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/SwaggerJwkUpdateSetKey} obj Optional instance to populate.
   * @return {module:model/SwaggerJwkUpdateSetKey} The populated <code>SwaggerJwkUpdateSetKey</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();

      if (data.hasOwnProperty('Body')) {
        obj['Body'] = JSONWebKey.constructFromObject(data['Body']);
      }
      if (data.hasOwnProperty('kid')) {
        obj['kid'] = ApiClient.convertToType(data['kid'], 'String');
      }
      if (data.hasOwnProperty('set')) {
        obj['set'] = ApiClient.convertToType(data['set'], 'String');
      }
    }
    return obj;
  }

  /**
   * @member {module:model/JSONWebKey} Body
   */
  exports.prototype['Body'] = undefined;
  /**
   * The kid of the desired key in: path
   * @member {String} kid
   */
  exports.prototype['kid'] = undefined;
  /**
   * The set in: path
   * @member {String} set
   */
  exports.prototype['set'] = undefined;



  return exports;
}));


