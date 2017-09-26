/**
 * Hydra OAuth2 & OpenID Connect Server (1.0.0-aplha1)
 * Please refer to the user guide for in-depth documentation: https://ory.gitbooks.io/hydra/content/   Hydra offers OAuth 2.0 and OpenID Connect Core 1.0 capabilities as a service. Hydra is different, because it works with any existing authentication infrastructure, not just LDAP or SAML. By implementing a consent app (works with any programming language) you build a bridge between Hydra and your authentication infrastructure. Hydra is able to securely manage JSON Web Keys, and has a sophisticated policy-based access control you can use if you want to. Hydra is suitable for green- (new) and brownfield (existing) projects. If you are not familiar with OAuth 2.0 and are working on a greenfield project, we recommend evaluating if OAuth 2.0 really serves your purpose. Knowledge of OAuth 2.0 is imperative in understanding what Hydra does and how it works.   The official repository is located at https://github.com/ory/hydra   ### ATTENTION - IMPORTANT NOTE   The swagger generator used to create this documentation does currently not support example responses. To see request and response payloads click on **\"Show JSON schema\"**: ![Enable JSON Schema on Apiary](https://storage.googleapis.com/ory.am/hydra/json-schema.png)
 *
 * OpenAPI spec version: Latest
 * Contact: hi@ory.am
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 2.2.3
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
    if (!root.HydraOAuth2OpenIdConnectServer100Aplha1) {
      root.HydraOAuth2OpenIdConnectServer100Aplha1 = {};
    }
    root.HydraOAuth2OpenIdConnectServer100Aplha1.Context = factory(root.HydraOAuth2OpenIdConnectServer100Aplha1.ApiClient);
  }
}(this, function(ApiClient) {
  'use strict';




  /**
   * The Context model module.
   * @module model/Context
   * @version Latest
   */

  /**
   * Constructs a new <code>Context</code>.
   * Context contains an access token&#39;s session data
   * @alias module:model/Context
   * @class
   */
  var exports = function() {
    var _this = this;






  };

  /**
   * Constructs a <code>Context</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/Context} obj Optional instance to populate.
   * @return {module:model/Context} The populated <code>Context</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();

      if (data.hasOwnProperty('aud')) {
        obj['aud'] = ApiClient.convertToType(data['aud'], 'String');
      }
      if (data.hasOwnProperty('ext')) {
        obj['ext'] = ApiClient.convertToType(data['ext'], {'String': Object});
      }
      if (data.hasOwnProperty('iss')) {
        obj['iss'] = ApiClient.convertToType(data['iss'], 'String');
      }
      if (data.hasOwnProperty('scopes')) {
        obj['scopes'] = ApiClient.convertToType(data['scopes'], ['String']);
      }
      if (data.hasOwnProperty('sub')) {
        obj['sub'] = ApiClient.convertToType(data['sub'], 'String');
      }
    }
    return obj;
  }

  /**
   * Audience is who the token was issued for. This is an OAuth2 app usually.
   * @member {String} aud
   */
  exports.prototype['aud'] = undefined;
  /**
   * Extra represents arbitrary session data.
   * @member {Object.<String, Object>} ext
   */
  exports.prototype['ext'] = undefined;
  /**
   * Issuer is the id of the issuer, typically an hydra instance.
   * @member {String} iss
   */
  exports.prototype['iss'] = undefined;
  /**
   * GrantedScopes is a list of scopes that the subject authorized when asked for consent.
   * @member {Array.<String>} scopes
   */
  exports.prototype['scopes'] = undefined;
  /**
   * Subject is the identity that authorized issuing the token, for example a user or an OAuth2 app. This is usually a uuid but you can choose a urn or some other id too.
   * @member {String} sub
   */
  exports.prototype['sub'] = undefined;



  return exports;
}));


