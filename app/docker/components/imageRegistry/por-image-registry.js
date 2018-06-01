angular.module('chainid.docker').component('porImageRegistry', {
  templateUrl: 'app/docker/components/imageRegistry/porImageRegistry.html',
  controller: 'porImageRegistryController',
  bindings: {
    'image': '=',
    'registry': '=',
    'autoComplete': '<'
  }
});
