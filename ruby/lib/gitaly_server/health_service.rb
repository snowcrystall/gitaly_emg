require 'grpc/health/v1/health_services_pb'

module GitalyServer
  class HealthService < Grpc::Health::V1::Health::Service
    def check(_req, _call)
      Grpc::Health::V1::HealthCheckResponse.new
    end
  end
end
