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
    root.OryHydra.RequestHandlerResponse = factory(root.OryHydra.ApiClient);
  }
}(this, function(ApiClient) {
  'use strict';




  /**
   * The RequestHandlerResponse model module.
   * @module model/RequestHandlerResponse
   * @version latest
   */

  /**
   * Constructs a new <code>RequestHandlerResponse</code>.
   * @alias module:model/RequestHandlerResponse
   * @class
   */
  var exports = function() {
    var _this = this;


  };

  /**
   * Constructs a <code>RequestHandlerResponse</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/RequestHandlerResponse} obj Optional instance to populate.
   * @return {module:model/RequestHandlerResponse} The populated <code>RequestHandlerResponse</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();

      if (data.hasOwnProperty('redirect_to')) {
        obj['redirect_to'] = ApiClient.convertToType(data['redirect_to'], 'String');
      }
    }
    return obj;
  }

  /**
   * RedirectURL is the URL which you should redirect the user to once the authentication process is completed.
   * @member {String} redirect_to
   */
  exports.prototype['redirect_to'] = undefined;



  return exports;
}));


