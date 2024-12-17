import Ember from 'ember';

export function permissionCheck(param) {
	if(param.length < 2){
		return false;
	}
	for(var i=0;i<param[0].length;i++){
		if(param[0][i] === param[1]){
			return true;
		}
	}
	return false;
}

export default Ember.Helper.helper(permissionCheck);