angular.module('chainid.docker')
.controller('ImageController', ['$q', '$scope', '$transition$', '$state', '$timeout', 'ImageService', 'RegistryService', 'Notifications', 'HttpRequestHelper',
function ($q, $scope, $transition$, $state, $timeout, ImageService, RegistryService, Notifications, HttpRequestHelper) {
	$scope.formValues = {
		Image: '',
		Registry: ''
	};

	$scope.sortType = 'Order';
	$scope.sortReverse = false;

	$scope.order = function(sortType) {
    $scope.sortReverse = ($scope.sortType === sortType) ? !$scope.sortReverse : false;
    $scope.sortType = sortType;
  };

	$scope.toggleLayerCommand = function(layerId) {
		$('#layer-command-expander'+layerId+' span').toggleClass('glyphicon-plus-sign glyphicon-minus-sign');
		$('#layer-command-'+layerId+'-short').toggle();
		$('#layer-command-'+layerId+'-full').toggle();
	};

	$scope.tagImage = function() {
		var image = $scope.formValues.Image;
		var registry = $scope.formValues.Registry;

		ImageService.tagImage($transition$.params().id, image, registry.URL)
		.then(function success(data) {
			Notifications.success('Image successfully tagged');
			$state.go('docker.images.image', {id: $transition$.params().id}, {reload: true});
		})
		.catch(function error(err) {
			Notifications.error('Failure', err, 'Unable to tag image');
		});
	};

	$scope.pushTag = function(repository) {
		$('#uploadResourceHint').show();
		RegistryService.retrieveRegistryFromRepository(repository)
		.then(function success(data) {
			var registry = data;
			return ImageService.pushImage(repository, registry);
		})
		.then(function success(data) {
			Notifications.success('Image successfully pushed', repository);
		})
		.catch(function error(err) {
			Notifications.error('Failure', err, 'Unable to push image to repository');
		})
		.finally(function final() {
			$('#uploadResourceHint').hide();
		});
	};

	$scope.pullTag = function(repository) {
		$('#downloadResourceHint').show();
		RegistryService.retrieveRegistryFromRepository(repository)
		.then(function success(data) {
			var registry = data;
			return ImageService.pullImage(repository, registry, false);
		})
		.then(function success(data) {
			Notifications.success('Image successfully pulled', repository);
		})
		.catch(function error(err) {
			Notifications.error('Failure', err, 'Unable to pull image');
		})
		.finally(function final() {
			$('#downloadResourceHint').hide();
		});
	};

	$scope.removeTag = function(repository) {
		ImageService.deleteImage(repository, false)
		.then(function success() {
			if ($scope.image.RepoTags.length === 1) {
				Notifications.success('Image successfully deleted', repository);
				$state.go('docker.images', {}, {reload: true});
			} else {
				Notifications.success('Tag successfully deleted', repository);
				$state.go('docker.images.image', {id: $transition$.params().id}, {reload: true});
			}
		})
		.catch(function error(err) {
			Notifications.error('Failure', err, 'Unable to remove image');
		});
	};

	$scope.removeImage = function (id) {
		ImageService.deleteImage(id, false)
		.then(function success() {
			Notifications.success('Image successfully deleted', id);
			$state.go('docker.images', {}, {reload: true});
		})
		.catch(function error(err) {
			Notifications.error('Failure', err, 'Unable to remove image');
		});
	};

	function initView() {
		HttpRequestHelper.setPortainerAgentTargetHeader($transition$.params().nodeName);
		var endpointProvider = $scope.applicationState.endpoint.mode.provider;
		$q.all({
			image: ImageService.image($transition$.params().id),
			history: endpointProvider !== 'VMWARE_VIC' ? ImageService.history($transition$.params().id) : []
		})
		.then(function success(data) {
			$scope.image = data.image;
			$scope.history = data.history;
		})
		.catch(function error(err) {
			Notifications.error('Failure', err, 'Unable to retrieve image details');
			$state.go('docker.images');
		});
	}

	initView();
}]);
