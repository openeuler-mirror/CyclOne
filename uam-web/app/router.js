import Ember from 'ember'
import config from './config/environment'

const Router = Ember.Router.extend({
    location: config.locationType
})

Router.map(function() {
    this.route('dashboard', function() {
        this.route('error')
        this.route('portalManagement', function() {
            this.route('userGroup', function() {
                this.route('operate')
                this.route('view', {
                    path: 'view/:id'
                })
            })
            this.route('user', function() {
              this.route('operate', {
                  path: 'operate/:id'
              })
              this.route('view', {
                  path: 'view/:id'
              })
              this.route('import');
              this.route('importPriview', {
                path: 'importPriview/:id'
              });
            })
            this.route('role', function() {
                this.route('operate')
                this.route('view', {
                    path: 'view/:id'
                })
            })
            this.route('permission', function() {
                this.route('reverse')
            })
            this.route('resource', function() {
                this.route('operate', {
                    path: 'operate/:id'
                })
                this.route('view', {
                    path: 'view/:id'
                })
            })
            this.route('tenant')
        })
    })
    this.route('authenticated')
    this.route('error', {
        path: 'error/:message'
    })
})

export default Router
