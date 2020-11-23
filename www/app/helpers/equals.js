import Ember from 'ember';

export default Ember.Helper.helper(function equals(params) {
    return params[0] === params[1];
});
