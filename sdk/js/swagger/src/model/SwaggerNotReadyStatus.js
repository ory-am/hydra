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
    define(['ApiClient'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'));
  } else {
    // Browser globals (root is window)
    if (!root.OryHydra) {
      root.OryHydra = {};
    }
    root.OryHydra.SwaggerNotReadyStatus = factory(root.OryHydra.ApiClient);
  }
}(this, function(ApiClient) {
  'use strict';




  /**
   * The SwaggerNotReadyStatus model module.
   * @module model/SwaggerNotReadyStatus
   * @version latest
   */

  /**
   * Constructs a new <code>SwaggerNotReadyStatus</code>.
   * SwaggerNotReadyStatus swagger not ready status
   * @alias module:model/SwaggerNotReadyStatus
   * @class
   */
  var exports = function() {
    var _this = this;


  };

  /**
   * Constructs a <code>SwaggerNotReadyStatus</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/SwaggerNotReadyStatus} obj Optional instance to populate.
   * @return {module:model/SwaggerNotReadyStatus} The populated <code>SwaggerNotReadyStatus</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();

      if (data.hasOwnProperty('errors')) {
        obj['errors'] = ApiClient.convertToType(data['errors'], {'String': 'String'});
      }
    }
    return obj;
  }

  /**
   * Errors contains a list of errors that caused the not ready status.
   * @member {Object.<String, String>} errors
   */
  exports.prototype['errors'] = undefined;



  return exports;
}));


