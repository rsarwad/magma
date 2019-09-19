/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @generated
 */

export type alert_bulk_upload_response = {
    errors: {
        [string]: string,
    },
    statuses: {
        [string]: string,
    },
};
export type alert_receiver_config = {
    name: string,
    slack_configs ? : Array < slack_receiver >
        ,
};
export type alert_routing_tree = {
    continue ?: boolean,
    group_by ? : Array < string >
        ,
    group_interval ? : string,
    group_wait ? : string,
    match ? : {
        label ? : string,
        value ? : string,
    },
    match_re ? : {
        label ? : string,
        value ? : string,
    },
    receiver: string,
    repeat_interval ? : string,
    routes ? : Array < alert_routing_tree >
        ,
};
export type challenge_key = {
    key ? : string,
    key_type: "ECHO" | "SOFTWARE_ECDSA_SHA256",
};
export type channel_id = string;
export type config_info = {
    mconfig_created_at ? : number,
};
export type disk_partition = {
    device ? : string,
    free ? : number,
    mount_point ? : string,
    total ? : number,
    used ? : number,
};
export type dns_config_record = {
    a_record ? : Array < string >
        ,
    aaaa_record ? : Array < string >
        ,
    cname_record ? : Array < string >
        ,
    domain: string,
};
export type enodeb = {
    attached_gateway_id ? : string,
    config: enodeb_configuration,
    name: string,
    serial: string,
};
export type enodeb_configuration = {
    bandwidth_mhz ? : 3 | 5 | 10 | 15 | 20,
    cell_id: number,
    device_class: "Baicells Nova-233 G2 OD FDD" | "Baicells Nova-243 OD TDD" | "Baicells Neutrino 224 ID FDD" | "Baicells ID TDD/FDD" | "NuRAN Cavium OC-LTE",
    earfcndl ? : number,
    pci ? : number,
    special_subframe_pattern ? : number,
    subframe_assignment ? : number,
    tac ? : number,
    transmit_enabled: boolean,
};
export type enodeb_serials = Array < string >
;
export type error = {
    message: string,
};
export type feg_network_id = string;
export type flow_record = {
    bytes_rx: number,
    bytes_tx: number,
    pkts_rx: number,
    pkts_tx: number,
    subscriber_id: string,
};
export type gateway_cellular_configs = {
    epc: gateway_epc_configs,
    non_eps_service ? : gateway_non_eps_configs,
    ran: gateway_ran_configs,
};
export type gateway_description = string;
export type gateway_device = {
    hardware_id: string,
    key: challenge_key,
};
export type gateway_epc_configs = {
    ip_block: string,
    nat_enabled: boolean,
};
export type gateway_id = string;
export type gateway_name = string;
export type gateway_non_eps_configs = {
    arfcn_2g: Array < number >
        ,
    csfb_mcc: string,
    csfb_mnc: string,
    csfb_rat: 0 | 1,
    lac: number,
    non_eps_service_control: 0 | 1 | 2,
};
export type gateway_ran_configs = {
    pci: number,
    transmit_enabled: boolean,
};
export type gateway_status = {
    cert_expiration_time ? : number,
    checkin_time ? : number,
    hardware_id ? : string,
    kernel_version ? : string,
    kernel_versions_installed ? : Array < string >
        ,
    machine_info ? : machine_info,
    meta ? : {
        [string]: string,
    },
    platform_info ? : platform_info,
    system_status ? : system_status,
    version ? : string,
    vpn_ip ? : string,
};
export type generic_command_params = {
    command: string,
    params ? : {
        [string]: {},
    },
};
export type generic_command_response = {
    response ? : {
        [string]: {},
    },
};
export type gettable_alert = {
    name: string,
};
export type label_pair = {
    name: string,
    value: string,
};
export type lte_gateway = {
    cellular: gateway_cellular_configs,
    connected_enodeb_serials: enodeb_serials,
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    status ? : gateway_status,
    tier: tier_id,
};
export type lte_network = {
    cellular: network_cellular_configs,
    description: network_description,
    dns: network_dns_config,
    features ? : network_features,
    id: network_id,
    name: network_name,
};
export type lte_subscription = {
    auth_algo: "MILENAGE",
    auth_key: string,
    auth_opc ? : string,
    state: "INACTIVE" | "ACTIVE",
};
export type machine_info = {
    cpu_info ? : {
        architecture ? : string,
        core_count ? : number,
        model_name ? : string,
        threads_per_core ? : number,
    },
    network_info ? : {
        network_interfaces ? : Array < network_interface >
            ,
        routing_table ? : Array < route >
            ,
    },
};
export type magmad_gateway = {
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    status ? : gateway_status,
    tier: tier_id,
};
export type magmad_gateway_configs = {
    autoupgrade_enabled: boolean,
    autoupgrade_poll_interval: number,
    checkin_interval: number,
    checkin_timeout: number,
    dynamic_services ? : Array < string >
        ,
    feature_flags ? : {
        [string]: boolean,
    },
};
export type metric_datapoint = Array < string >
;
export type metric_datapoints = Array < metric_datapoint >
;
export type mutable_lte_gateway = {
    cellular: gateway_cellular_configs,
    connected_enodeb_serials: enodeb_serials,
    description: gateway_description,
    device: gateway_device,
    id: gateway_id,
    magmad: magmad_gateway_configs,
    name: gateway_name,
    tier: tier_id,
};
export type network = {
    description: network_description,
    dns: network_dns_config,
    features ? : network_features,
    id: network_id,
    name: network_name,
    type ? : network_type,
};
export type network_cellular_configs = {
    epc: network_epc_configs,
    feg_network_id ? : feg_network_id,
    ran: network_ran_configs,
};
export type network_description = string;
export type network_dns_config = {
    enable_caching: boolean,
    local_ttl: number,
    records ? : network_dns_records,
};
export type network_dns_records = Array < dns_config_record >
;
export type network_epc_configs = {
    cloud_subscriberdb_enabled ? : boolean,
    default_rule_id ? : string,
    lte_auth_amf: string,
    lte_auth_op: string,
    mcc: string,
    mnc: string,
    mobility ? : {
        ip_allocation_mode: "NAT" | "STATIC" | "DHCP_PASSTHROUGH" | "DHCP_BROADCAST",
        nat ? : {
            ip_blocks ? : Array < string >
                ,
        },
        reserved_addresses ? : Array < string >
            ,
        static ? : {
            ip_blocks_by_tac ? : {
                [string]: Array < string >
                    ,
            },
        },
    },
    network_services ? : Array < "metering" | "dpi" | "policy_enforcement" >
        ,
    relay_enabled: boolean,
    sub_profiles ? : {
        [string]: {
            max_dl_bit_rate: number,
            max_ul_bit_rate: number,
        },
    },
    tac: number,
};
export type network_features = {
    features ? : {
        [string]: string,
    },
};
export type network_id = string;
export type network_interface = {
    ip_addresses ? : Array < string >
        ,
    ipv6_addresses ? : Array < string >
        ,
    mac_address ? : string,
    network_interface_id ? : string,
    status ? : "UP" | "DOWN" | "UNKNOWN",
};
export type network_name = string;
export type network_ran_configs = {
    bandwidth_mhz: 3 | 5 | 10 | 15 | 20,
    fdd_config ? : {
        earfcndl: number,
        earfcnul: number,
    },
    tdd_config ? : {
        earfcndl: number,
        special_subframe_pattern: number,
        subframe_assignment: number,
    },
};
export type network_type = string;
export type package_type = {
    name ? : string,
    version ? : string,
};
export type ping_request = {
    hosts: Array < string >
        ,
    packets ? : number,
};
export type ping_response = {
    pings: Array < ping_result >
        ,
};
export type ping_result = {
    avg_response_ms ? : number,
    error ? : string,
    host_or_ip: string,
    num_packets: number,
    packets_received ? : number,
    packets_transmitted ? : number,
};
export type platform_info = {
    config_info ? : config_info,
    kernel_version ? : string,
    kernel_versions_installed ? : Array < string >
        ,
    packages ? : Array < package_type >
        ,
    vpn_ip ? : string,
};
export type prom_alert_config = {
    alert: string,
    annotations ? : prom_alert_labels,
    expr: string,
    for ? : string,
    labels ? : prom_alert_labels,
};
export type prom_alert_config_list = Array < prom_alert_config >
;
export type prom_alert_labels = {
    [string]: string,
};
export type prom_alert_status = {
    inhibitedBy: Array < string >
        ,
    silencedBy: Array < string >
        ,
    state: string,
};
export type prom_firing_alert = {
    annotations: prom_alert_labels,
    endsAt: string,
    fingerprint: string,
    generatorURL ? : string,
    labels: prom_alert_labels,
    receivers: gettable_alert,
    startsAt: string,
    status: prom_alert_status,
    updatedAt: string,
};
export type promql_data = {
    result: promql_result,
    resultType: string,
};
export type promql_metric = {
    __name__: string,
    gateway ? : string,
    host ? : string,
    instance: string,
    job ? : string,
};
export type promql_metric_value = {
    metric: promql_metric,
    value ? : metric_datapoint,
    values ? : metric_datapoints,
};
export type promql_result = Array < promql_metric_value >
;
export type promql_return_object = {
    data: promql_data,
    status: string,
};
export type pushed_metric = {
    labels ? : Array < label_pair >
        ,
    metricName: string,
    timestamp ? : string,
    value: number,
};
export type release_channel = {
    id: channel_id,
    name ? : string,
    supported_versions: Array < string >
        ,
};
export type route = {
    destination_ip ? : string,
    gateway_ip ? : string,
    genmask ? : string,
    network_interface_id ? : string,
};
export type slack_action = {
    confirm ? : slack_confirm_field,
    name ? : string,
    style ? : string,
    text: string,
    type: string,
    url: string,
    value ? : string,
};
export type slack_confirm_field = {
    dismiss_text: string,
    ok_text: string,
    text: string,
    title: string,
};
export type slack_field = {
    short ? : boolean,
    title: string,
    value: string,
};
export type slack_receiver = {
    actions ? : Array < slack_action >
        ,
    api_url: string,
    callback_id ? : string,
    channel ? : string,
    color ? : string,
    fallback ? : string,
    fields ? : Array < slack_field >
        ,
    footer ? : string,
    icon_emoji ? : string,
    icon_url ? : string,
    image_url ? : string,
    link_names ? : boolean,
    pretext ? : string,
    short_fields ? : boolean,
    text ? : string,
    thumb_url ? : string,
    title ? : string,
    username ? : string,
};
export type subscriber = {
    id: string,
    lte: lte_subscription,
};
export type symphony_network = {
    description: network_description,
    features ? : network_features,
    id: network_id,
    name: network_name,
};
export type system_status = {
    cpu_idle ? : number,
    cpu_system ? : number,
    cpu_user ? : number,
    disk_partitions ? : Array < disk_partition >
        ,
    mem_available ? : number,
    mem_free ? : number,
    mem_total ? : number,
    mem_used ? : number,
    swap_free ? : number,
    swap_total ? : number,
    swap_used ? : number,
    time ? : number,
    uptime_secs ? : number,
};
export type tail_logs_request = {
    service ? : string,
};
export type tier = {
    gateways: tier_gateways,
    id: tier_id,
    images: tier_images,
    name ? : tier_name,
    version: tier_version,
};
export type tier_gateways = Array < gateway_id >
;
export type tier_id = string;
export type tier_image = {
    name: string,
    order: number,
};
export type tier_images = Array < tier_image >
;
export type tier_name = string;
export type tier_version = string;

export default class MagmaAPIBindings {
    static request(
        path: string,
        method: 'POST' | 'GET' | 'PUT' | 'DELETE' | 'OPTIONS' | 'HEAD' | 'PATCH',
        query: {
            [string]: mixed
        },
        body ? : {
            [string]: any
        } | string | Array < any > ,
    ) {
        throw new Error("Must be implemented");
    }
    static async getChannels(): Promise < Array < channel_id >
        >
        {
            let path = '/channels';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postChannels(
        parameters: {
            'channel': release_channel,
        }
    ): Promise < "Success" > {
        let path = '/channels';
        let body;
        let query = {};
        if (parameters['channel'] === undefined) {
            throw new Error('Missing required  parameter: channel');
        }

        if (parameters['channel'] !== undefined) {
            body = parameters['channel'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteChannelsByChannelId(
        parameters: {
            'channelId': string,
        }
    ): Promise < "Success" > {
        let path = '/channels/{channel_id}';
        let body;
        let query = {};
        if (parameters['channelId'] === undefined) {
            throw new Error('Missing required  parameter: channelId');
        }

        path = path.replace('{channel_id}', `${parameters['channelId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getChannelsByChannelId(
            parameters: {
                'channelId': string,
            }
        ): Promise < release_channel >
        {
            let path = '/channels/{channel_id}';
            let body;
            let query = {};
            if (parameters['channelId'] === undefined) {
                throw new Error('Missing required  parameter: channelId');
            }

            path = path.replace('{channel_id}', `${parameters['channelId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putChannelsByChannelId(
        parameters: {
            'channelId': string,
            'releaseChannel': release_channel,
        }
    ): Promise < "Success" > {
        let path = '/channels/{channel_id}';
        let body;
        let query = {};
        if (parameters['channelId'] === undefined) {
            throw new Error('Missing required  parameter: channelId');
        }

        path = path.replace('{channel_id}', `${parameters['channelId']}`);

        if (parameters['releaseChannel'] === undefined) {
            throw new Error('Missing required  parameter: releaseChannel');
        }

        if (parameters['releaseChannel'] !== undefined) {
            body = parameters['releaseChannel'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getFoo(): Promise < number >
        {
            let path = '/foo';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async getLte(): Promise < Array < string >
        >
        {
            let path = '/lte';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postLte(
        parameters: {
            'lteNetwork': lte_network,
        }
    ): Promise < "Success" > {
        let path = '/lte';
        let body;
        let query = {};
        if (parameters['lteNetwork'] === undefined) {
            throw new Error('Missing required  parameter: lteNetwork');
        }

        if (parameters['lteNetwork'] !== undefined) {
            body = parameters['lteNetwork'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteLteByNetworkId(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < lte_network >
        {
            let path = '/lte/{network_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkId(
        parameters: {
            'networkId': string,
            'lteNetwork': lte_network,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['lteNetwork'] === undefined) {
            throw new Error('Missing required  parameter: lteNetwork');
        }

        if (parameters['lteNetwork'] !== undefined) {
            body = parameters['lteNetwork'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdCellular(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_cellular_configs >
        {
            let path = '/lte/{network_id}/cellular';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdCellular(
        parameters: {
            'networkId': string,
            'config': network_cellular_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/cellular';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdCellularEpc(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_epc_configs >
        {
            let path = '/lte/{network_id}/cellular/epc';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdCellularEpc(
        parameters: {
            'networkId': string,
            'config': network_epc_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/cellular/epc';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdCellularFegNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < string >
        {
            let path = '/lte/{network_id}/cellular/feg_network_id';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdCellularFegNetworkId(
        parameters: {
            'networkId': string,
            'fegNetworkId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/cellular/feg_network_id';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['fegNetworkId'] === undefined) {
            throw new Error('Missing required  parameter: fegNetworkId');
        }

        if (parameters['fegNetworkId'] !== undefined) {
            body = parameters['fegNetworkId'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdCellularRan(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_ran_configs >
        {
            let path = '/lte/{network_id}/cellular/ran';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdCellularRan(
        parameters: {
            'networkId': string,
            'config': network_ran_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/cellular/ran';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdDescription(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_description >
        {
            let path = '/lte/{network_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdDescription(
        parameters: {
            'networkId': string,
            'description': network_description,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdDns(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_dns_config >
        {
            let path = '/lte/{network_id}/dns';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdDns(
        parameters: {
            'networkId': string,
            'config': network_dns_config,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/dns';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdDnsRecords(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < dns_config_record >
        >
        {
            let path = '/lte/{network_id}/dns/records';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdDnsRecords(
        parameters: {
            'networkId': string,
            'records': Array < dns_config_record >
                ,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/dns/records';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['records'] === undefined) {
            throw new Error('Missing required  parameter: records');
        }

        if (parameters['records'] !== undefined) {
            body = parameters['records'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteLteByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkIdDnsRecordsByDomain(
            parameters: {
                'networkId': string,
                'domain': string,
            }
        ): Promise < dns_config_record >
        {
            let path = '/lte/{network_id}/dns/records/{domain}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['domain'] === undefined) {
                throw new Error('Missing required  parameter: domain');
            }

            path = path.replace('{domain}', `${parameters['domain']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postLteByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
            'record': dns_config_record,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putLteByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
            'record': dns_config_record,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdEnodebs(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: enodeb,
        } >
        {
            let path = '/lte/{network_id}/enodebs';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postLteByNetworkIdEnodebs(
        parameters: {
            'networkId': string,
            'enodeb': enodeb,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/enodebs';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['enodeb'] === undefined) {
            throw new Error('Missing required  parameter: enodeb');
        }

        if (parameters['enodeb'] !== undefined) {
            body = parameters['enodeb'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteLteByNetworkIdEnodebsByEnodebSerial(
        parameters: {
            'networkId': string,
            'enodebSerial': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/enodebs/{enodeb_serial}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['enodebSerial'] === undefined) {
            throw new Error('Missing required  parameter: enodebSerial');
        }

        path = path.replace('{enodeb_serial}', `${parameters['enodebSerial']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkIdEnodebsByEnodebSerial(
            parameters: {
                'networkId': string,
                'enodebSerial': string,
            }
        ): Promise < enodeb >
        {
            let path = '/lte/{network_id}/enodebs/{enodeb_serial}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['enodebSerial'] === undefined) {
                throw new Error('Missing required  parameter: enodebSerial');
            }

            path = path.replace('{enodeb_serial}', `${parameters['enodebSerial']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdEnodebsByEnodebSerial(
        parameters: {
            'networkId': string,
            'enodebSerial': string,
            'enodeb': enodeb,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/enodebs/{enodeb_serial}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['enodebSerial'] === undefined) {
            throw new Error('Missing required  parameter: enodebSerial');
        }

        path = path.replace('{enodeb_serial}', `${parameters['enodebSerial']}`);

        if (parameters['enodeb'] === undefined) {
            throw new Error('Missing required  parameter: enodeb');
        }

        if (parameters['enodeb'] !== undefined) {
            body = parameters['enodeb'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdFeatures(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_features >
        {
            let path = '/lte/{network_id}/features';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdFeatures(
        parameters: {
            'networkId': string,
            'config': network_features,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/features';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGateways(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: lte_gateway,
        } >
        {
            let path = '/lte/{network_id}/gateways';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postLteByNetworkIdGateways(
        parameters: {
            'networkId': string,
            'gateway': mutable_lte_gateway,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteLteByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayId(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < lte_gateway >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'gateway': mutable_lte_gateway,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdCellular(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_cellular_configs >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/cellular';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdCellular(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_cellular_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/cellular';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdCellularEpc(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_epc_configs >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/epc';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdCellularEpc(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_epc_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/epc';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdCellularNonEps(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_non_eps_configs >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/non_eps';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdCellularNonEps(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_non_eps_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/non_eps';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdCellularRan(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_ran_configs >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/ran';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdCellularRan(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'config': gateway_ran_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/cellular/ran';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteLteByNetworkIdGatewaysByGatewayIdConnectedEnodebSerials(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'serial': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/connected_enodeb_serials';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['serial'] === undefined) {
            throw new Error('Missing required  parameter: serial');
        }

        if (parameters['serial'] !== undefined) {
            body = parameters['serial'];
        }

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdConnectedEnodebSerials(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < enodeb_serials >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/connected_enodeb_serials';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postLteByNetworkIdGatewaysByGatewayIdConnectedEnodebSerials(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'serial': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/connected_enodeb_serials';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['serial'] === undefined) {
            throw new Error('Missing required  parameter: serial');
        }

        if (parameters['serial'] !== undefined) {
            body = parameters['serial'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putLteByNetworkIdGatewaysByGatewayIdConnectedEnodebSerials(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'serials': enodeb_serials,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/connected_enodeb_serials';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['serials'] === undefined) {
            throw new Error('Missing required  parameter: serials');
        }

        if (parameters['serials'] !== undefined) {
            body = parameters['serials'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdDescription(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_description >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdDescription(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'description': gateway_description,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdDevice(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_device >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/device';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdDevice(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'device': gateway_device,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/device';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['device'] === undefined) {
            throw new Error('Missing required  parameter: device');
        }

        if (parameters['device'] !== undefined) {
            body = parameters['device'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdMagmad(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < magmad_gateway_configs >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/magmad';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdMagmad(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'magmad': magmad_gateway_configs,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/magmad';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['magmad'] === undefined) {
            throw new Error('Missing required  parameter: magmad');
        }

        if (parameters['magmad'] !== undefined) {
            body = parameters['magmad'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdName(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_name >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdName(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'name': gateway_name,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdGatewaysByGatewayIdStatus(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_status >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/status';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getLteByNetworkIdGatewaysByGatewayIdTier(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < tier_id >
        {
            let path = '/lte/{network_id}/gateways/{gateway_id}/tier';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdGatewaysByGatewayIdTier(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'tierId': tier_id,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/gateways/{gateway_id}/tier';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        if (parameters['tierId'] !== undefined) {
            body = parameters['tierId'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdName(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_name >
        {
            let path = '/lte/{network_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdName(
        parameters: {
            'networkId': string,
            'name': network_name,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getLteByNetworkIdSubscribers(
            parameters: {
                'networkId': string,
            }
        ): Promise < {
            [string]: subscriber,
        } >
        {
            let path = '/lte/{network_id}/subscribers';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postLteByNetworkIdSubscribers(
        parameters: {
            'networkId': string,
            'subscriber': subscriber,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriber'] === undefined) {
            throw new Error('Missing required  parameter: subscriber');
        }

        if (parameters['subscriber'] !== undefined) {
            body = parameters['subscriber'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteLteByNetworkIdSubscribersBySubscriberId(
        parameters: {
            'networkId': string,
            'subscriberId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers/{subscriber_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriberId'] === undefined) {
            throw new Error('Missing required  parameter: subscriberId');
        }

        path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getLteByNetworkIdSubscribersBySubscriberId(
            parameters: {
                'networkId': string,
                'subscriberId': string,
            }
        ): Promise < subscriber >
        {
            let path = '/lte/{network_id}/subscribers/{subscriber_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['subscriberId'] === undefined) {
                throw new Error('Missing required  parameter: subscriberId');
            }

            path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putLteByNetworkIdSubscribersBySubscriberId(
        parameters: {
            'networkId': string,
            'subscriberId': string,
            'subscriber': subscriber,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers/{subscriber_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriberId'] === undefined) {
            throw new Error('Missing required  parameter: subscriberId');
        }

        path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

        if (parameters['subscriber'] === undefined) {
            throw new Error('Missing required  parameter: subscriber');
        }

        if (parameters['subscriber'] !== undefined) {
            body = parameters['subscriber'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async postLteByNetworkIdSubscribersBySubscriberIdActivate(
        parameters: {
            'networkId': string,
            'subscriberId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers/{subscriber_id}/activate';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriberId'] === undefined) {
            throw new Error('Missing required  parameter: subscriberId');
        }

        path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async postLteByNetworkIdSubscribersBySubscriberIdDeactivate(
        parameters: {
            'networkId': string,
            'subscriberId': string,
        }
    ): Promise < "Success" > {
        let path = '/lte/{network_id}/subscribers/{subscriber_id}/deactivate';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['subscriberId'] === undefined) {
            throw new Error('Missing required  parameter: subscriberId');
        }

        path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async getNetworks(): Promise < Array < string >
        >
        {
            let path = '/networks';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworks(
        parameters: {
            'network': network,
        }
    ): Promise < "Success" > {
        let path = '/networks';
        let body;
        let query = {};
        if (parameters['network'] === undefined) {
            throw new Error('Missing required  parameter: network');
        }

        if (parameters['network'] !== undefined) {
            body = parameters['network'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteNetworksByNetworkId(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < network >
        {
            let path = '/networks/{network_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkId(
        parameters: {
            'networkId': string,
            'network': network,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['network'] === undefined) {
            throw new Error('Missing required  parameter: network');
        }

        if (parameters['network'] !== undefined) {
            body = parameters['network'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdAlerts(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < prom_firing_alert >
        >
        {
            let path = '/networks/{network_id}/alerts';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getNetworksByNetworkIdDescription(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_description >
        {
            let path = '/networks/{network_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdDescription(
        parameters: {
            'networkId': string,
            'description': network_description,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdDns(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_dns_config >
        {
            let path = '/networks/{network_id}/dns';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdDns(
        parameters: {
            'networkId': string,
            'networkDns': network_dns_config,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/dns';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['networkDns'] === undefined) {
            throw new Error('Missing required  parameter: networkDns');
        }

        if (parameters['networkDns'] !== undefined) {
            body = parameters['networkDns'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdDnsRecords(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_dns_records >
        {
            let path = '/networks/{network_id}/dns/records';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdDnsRecords(
        parameters: {
            'networkId': string,
            'records': network_dns_records,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/dns/records';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['records'] === undefined) {
            throw new Error('Missing required  parameter: records');
        }

        if (parameters['records'] !== undefined) {
            body = parameters['records'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteNetworksByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdDnsRecordsByDomain(
            parameters: {
                'networkId': string,
                'domain': string,
            }
        ): Promise < dns_config_record >
        {
            let path = '/networks/{network_id}/dns/records/{domain}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['domain'] === undefined) {
                throw new Error('Missing required  parameter: domain');
            }

            path = path.replace('{domain}', `${parameters['domain']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
            'record': dns_config_record,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putNetworksByNetworkIdDnsRecordsByDomain(
        parameters: {
            'networkId': string,
            'domain': string,
            'record': dns_config_record,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/dns/records/{domain}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['domain'] === undefined) {
            throw new Error('Missing required  parameter: domain');
        }

        path = path.replace('{domain}', `${parameters['domain']}`);

        if (parameters['record'] === undefined) {
            throw new Error('Missing required  parameter: record');
        }

        if (parameters['record'] !== undefined) {
            body = parameters['record'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdFeatures(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_features >
        {
            let path = '/networks/{network_id}/features';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdFeatures(
        parameters: {
            'networkId': string,
            'networkFeatures': network_features,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/features';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['networkFeatures'] === undefined) {
            throw new Error('Missing required  parameter: networkFeatures');
        }

        if (parameters['networkFeatures'] !== undefined) {
            body = parameters['networkFeatures'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdGateways(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < magmad_gateway >
        >
        {
            let path = '/networks/{network_id}/gateways';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdGateways(
        parameters: {
            'networkId': string,
            'gateway': magmad_gateway,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteNetworksByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayId(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < magmad_gateway >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'gateway': magmad_gateway,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async postNetworksByNetworkIdGatewaysByGatewayIdCommandGeneric(
            parameters: {
                'networkId': string,
                'gatewayId': string,
                'parameters': generic_command_params,
            }
        ): Promise < generic_command_response >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/command/generic';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            if (parameters['parameters'] === undefined) {
                throw new Error('Missing required  parameter: parameters');
            }

            if (parameters['parameters'] !== undefined) {
                body = parameters['parameters'];
            }

            return await this.request(path, 'POST', query, body);
        }
    static async postNetworksByNetworkIdGatewaysByGatewayIdCommandPing(
            parameters: {
                'networkId': string,
                'gatewayId': string,
                'pingRequest': ping_request,
            }
        ): Promise < ping_response >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/command/ping';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            if (parameters['pingRequest'] === undefined) {
                throw new Error('Missing required  parameter: pingRequest');
            }

            if (parameters['pingRequest'] !== undefined) {
                body = parameters['pingRequest'];
            }

            return await this.request(path, 'POST', query, body);
        }
    static async postNetworksByNetworkIdGatewaysByGatewayIdCommandReboot(
        parameters: {
            'networkId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/command/reboot';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'POST', query, body);
    }
    static async postNetworksByNetworkIdGatewaysByGatewayIdCommandRestartServices(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'services': Array < string >
                ,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/command/restart_services';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['services'] === undefined) {
            throw new Error('Missing required  parameter: services');
        }

        if (parameters['services'] !== undefined) {
            body = parameters['services'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayIdDescription(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_description >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/description';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayIdDescription(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'description': gateway_description,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/description';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['description'] === undefined) {
            throw new Error('Missing required  parameter: description');
        }

        if (parameters['description'] !== undefined) {
            body = parameters['description'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayIdDevice(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_device >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/device';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayIdDevice(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'device': gateway_device,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/device';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['device'] === undefined) {
            throw new Error('Missing required  parameter: device');
        }

        if (parameters['device'] !== undefined) {
            body = parameters['device'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayIdMagmad(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < magmad_gateway_configs >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/magmad';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayIdMagmad(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'magmad': magmad_gateway_configs,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/magmad';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['magmad'] === undefined) {
            throw new Error('Missing required  parameter: magmad');
        }

        if (parameters['magmad'] !== undefined) {
            body = parameters['magmad'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayIdName(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_name >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayIdName(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'name': gateway_name,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdGatewaysByGatewayIdStatus(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < gateway_status >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/status';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getNetworksByNetworkIdGatewaysByGatewayIdTier(
            parameters: {
                'networkId': string,
                'gatewayId': string,
            }
        ): Promise < tier_id >
        {
            let path = '/networks/{network_id}/gateways/{gateway_id}/tier';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['gatewayId'] === undefined) {
                throw new Error('Missing required  parameter: gatewayId');
            }

            path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdGatewaysByGatewayIdTier(
        parameters: {
            'networkId': string,
            'gatewayId': string,
            'tierId': tier_id,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/gateways/{gateway_id}/tier';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        if (parameters['tierId'] !== undefined) {
            body = parameters['tierId'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async postNetworksByNetworkIdMetricsPush(
        parameters: {
            'networkId': string,
            'metrics': Array < pushed_metric >
                ,
        }
    ): Promise < "Submitted" > {
        let path = '/networks/{network_id}/metrics/push';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['metrics'] === undefined) {
            throw new Error('Missing required  parameter: metrics');
        }

        if (parameters['metrics'] !== undefined) {
            body = parameters['metrics'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async getNetworksByNetworkIdName(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_name >
        {
            let path = '/networks/{network_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdName(
        parameters: {
            'networkId': string,
            'name': network_name,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteNetworksByNetworkIdPrometheusAlertConfig(
        parameters: {
            'networkId': string,
            'alertName': string,
        }
    ): Promise < "Deleted" > {
        let path = '/networks/{network_id}/prometheus/alert_config';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['alertName'] === undefined) {
            throw new Error('Missing required  parameter: alertName');
        }

        if (parameters['alertName'] !== undefined) {
            query['alert_name'] = parameters['alertName'];
        }

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdPrometheusAlertConfig(
            parameters: {
                'networkId': string,
                'alertName' ? : string,
            }
        ): Promise < prom_alert_config_list >
        {
            let path = '/networks/{network_id}/prometheus/alert_config';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['alertName'] !== undefined) {
                query['alert_name'] = parameters['alertName'];
            }

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdPrometheusAlertConfig(
        parameters: {
            'networkId': string,
            'alertConfig': prom_alert_config,
        }
    ): Promise < "Created" > {
        let path = '/networks/{network_id}/prometheus/alert_config';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['alertConfig'] === undefined) {
            throw new Error('Missing required  parameter: alertConfig');
        }

        if (parameters['alertConfig'] !== undefined) {
            body = parameters['alertConfig'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putNetworksByNetworkIdPrometheusAlertConfigByAlertName(
        parameters: {
            'networkId': string,
            'alertName': string,
            'alertConfig': prom_alert_config,
        }
    ): Promise < "Updated" > {
        let path = '/networks/{network_id}/prometheus/alert_config/{alert_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['alertName'] === undefined) {
            throw new Error('Missing required  parameter: alertName');
        }

        path = path.replace('{alert_name}', `${parameters['alertName']}`);

        if (parameters['alertConfig'] === undefined) {
            throw new Error('Missing required  parameter: alertConfig');
        }

        if (parameters['alertConfig'] !== undefined) {
            body = parameters['alertConfig'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async putNetworksByNetworkIdPrometheusAlertConfigBulk(
            parameters: {
                'networkId': string,
                'alertConfigs': prom_alert_config_list,
            }
        ): Promise < alert_bulk_upload_response >
        {
            let path = '/networks/{network_id}/prometheus/alert_config/bulk';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['alertConfigs'] === undefined) {
                throw new Error('Missing required  parameter: alertConfigs');
            }

            if (parameters['alertConfigs'] !== undefined) {
                body = parameters['alertConfigs'];
            }

            return await this.request(path, 'PUT', query, body);
        }
    static async deleteNetworksByNetworkIdPrometheusAlertReceiver(
        parameters: {
            'networkId': string,
            'receiver': string,
        }
    ): Promise < "Deleted" > {
        let path = '/networks/{network_id}/prometheus/alert_receiver';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['receiver'] === undefined) {
            throw new Error('Missing required  parameter: receiver');
        }

        if (parameters['receiver'] !== undefined) {
            query['receiver'] = parameters['receiver'];
        }

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdPrometheusAlertReceiver(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < alert_receiver_config >
        >
        {
            let path = '/networks/{network_id}/prometheus/alert_receiver';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdPrometheusAlertReceiver(
        parameters: {
            'networkId': string,
            'receiverConfig': alert_receiver_config,
        }
    ): Promise < "Created" > {
        let path = '/networks/{network_id}/prometheus/alert_receiver';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['receiverConfig'] === undefined) {
            throw new Error('Missing required  parameter: receiverConfig');
        }

        if (parameters['receiverConfig'] !== undefined) {
            body = parameters['receiverConfig'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putNetworksByNetworkIdPrometheusAlertReceiverByReceiver(
        parameters: {
            'networkId': string,
            'receiver': string,
            'receiverConfig': alert_receiver_config,
        }
    ): Promise < "Updated" > {
        let path = '/networks/{network_id}/prometheus/alert_receiver/{receiver}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['receiver'] === undefined) {
            throw new Error('Missing required  parameter: receiver');
        }

        path = path.replace('{receiver}', `${parameters['receiver']}`);

        if (parameters['receiverConfig'] === undefined) {
            throw new Error('Missing required  parameter: receiverConfig');
        }

        if (parameters['receiverConfig'] !== undefined) {
            body = parameters['receiverConfig'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdPrometheusAlertReceiverRoute(
            parameters: {
                'networkId': string,
            }
        ): Promise < alert_routing_tree >
        {
            let path = '/networks/{network_id}/prometheus/alert_receiver/route';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdPrometheusAlertReceiverRoute(
        parameters: {
            'networkId': string,
            'route': alert_routing_tree,
        }
    ): Promise < "OK" > {
        let path = '/networks/{network_id}/prometheus/alert_receiver/route';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['route'] === undefined) {
            throw new Error('Missing required  parameter: route');
        }

        if (parameters['route'] !== undefined) {
            body = parameters['route'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async getNetworksByNetworkIdPrometheusQuery(
            parameters: {
                'networkId': string,
                'query': string,
                'time' ? : string,
            }
        ): Promise < Array < promql_return_object >
        >
        {
            let path = '/networks/{network_id}/prometheus/query';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['query'] === undefined) {
                throw new Error('Missing required  parameter: query');
            }

            if (parameters['query'] !== undefined) {
                query['query'] = parameters['query'];
            }

            if (parameters['time'] !== undefined) {
                query['time'] = parameters['time'];
            }

            return await this.request(path, 'GET', query, body);
        }
    static async getNetworksByNetworkIdPrometheusQueryRange(
            parameters: {
                'networkId': string,
                'query': string,
                'start': string,
                'end' ? : string,
                'step' ? : string,
            }
        ): Promise < Array < promql_return_object >
        >
        {
            let path = '/networks/{network_id}/prometheus/query_range';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['query'] === undefined) {
                throw new Error('Missing required  parameter: query');
            }

            if (parameters['query'] !== undefined) {
                query['query'] = parameters['query'];
            }

            if (parameters['start'] === undefined) {
                throw new Error('Missing required  parameter: start');
            }

            if (parameters['start'] !== undefined) {
                query['start'] = parameters['start'];
            }

            if (parameters['end'] !== undefined) {
                query['end'] = parameters['end'];
            }

            if (parameters['step'] !== undefined) {
                query['step'] = parameters['step'];
            }

            return await this.request(path, 'GET', query, body);
        }
    static async getNetworksByNetworkIdSubscribersBySubscriberIdFlowRecords(
            parameters: {
                'networkId': string,
                'subscriberId': string,
            }
        ): Promise < Array < flow_record >
        >
        {
            let path = '/networks/{network_id}/subscribers/{subscriber_id}/flow_records';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['subscriberId'] === undefined) {
                throw new Error('Missing required  parameter: subscriberId');
            }

            path = path.replace('{subscriber_id}', `${parameters['subscriberId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async getNetworksByNetworkIdTiers(
            parameters: {
                'networkId': string,
            }
        ): Promise < Array < tier_id >
        >
        {
            let path = '/networks/{network_id}/tiers';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdTiers(
        parameters: {
            'networkId': string,
            'tier': tier,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tier'] === undefined) {
            throw new Error('Missing required  parameter: tier');
        }

        if (parameters['tier'] !== undefined) {
            body = parameters['tier'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteNetworksByNetworkIdTiersByTierId(
        parameters: {
            'networkId': string,
            'tierId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdTiersByTierId(
            parameters: {
                'networkId': string,
                'tierId': string,
            }
        ): Promise < tier >
        {
            let path = '/networks/{network_id}/tiers/{tier_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['tierId'] === undefined) {
                throw new Error('Missing required  parameter: tierId');
            }

            path = path.replace('{tier_id}', `${parameters['tierId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdTiersByTierId(
        parameters: {
            'networkId': string,
            'tierId': string,
            'tier': tier,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['tier'] === undefined) {
            throw new Error('Missing required  parameter: tier');
        }

        if (parameters['tier'] !== undefined) {
            body = parameters['tier'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdTiersByTierIdGateways(
            parameters: {
                'networkId': string,
                'tierId': string,
            }
        ): Promise < tier_gateways >
        {
            let path = '/networks/{network_id}/tiers/{tier_id}/gateways';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['tierId'] === undefined) {
                throw new Error('Missing required  parameter: tierId');
            }

            path = path.replace('{tier_id}', `${parameters['tierId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdTiersByTierIdGateways(
        parameters: {
            'networkId': string,
            'tierId': string,
            'gateway': gateway_id,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['gateway'] === undefined) {
            throw new Error('Missing required  parameter: gateway');
        }

        if (parameters['gateway'] !== undefined) {
            body = parameters['gateway'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putNetworksByNetworkIdTiersByTierIdGateways(
        parameters: {
            'networkId': string,
            'tierId': string,
            'tier': tier_gateways,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/gateways';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['tier'] === undefined) {
            throw new Error('Missing required  parameter: tier');
        }

        if (parameters['tier'] !== undefined) {
            body = parameters['tier'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteNetworksByNetworkIdTiersByTierIdGatewaysByGatewayId(
        parameters: {
            'networkId': string,
            'tierId': string,
            'gatewayId': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/gateways/{gateway_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['gatewayId'] === undefined) {
            throw new Error('Missing required  parameter: gatewayId');
        }

        path = path.replace('{gateway_id}', `${parameters['gatewayId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdTiersByTierIdImages(
            parameters: {
                'networkId': string,
                'tierId': string,
            }
        ): Promise < tier_images >
        {
            let path = '/networks/{network_id}/tiers/{tier_id}/images';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['tierId'] === undefined) {
                throw new Error('Missing required  parameter: tierId');
            }

            path = path.replace('{tier_id}', `${parameters['tierId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async postNetworksByNetworkIdTiersByTierIdImages(
        parameters: {
            'networkId': string,
            'tierId': string,
            'image': tier_image,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/images';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['image'] === undefined) {
            throw new Error('Missing required  parameter: image');
        }

        if (parameters['image'] !== undefined) {
            body = parameters['image'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async putNetworksByNetworkIdTiersByTierIdImages(
        parameters: {
            'networkId': string,
            'tierId': string,
            'tier': tier_images,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/images';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['tier'] === undefined) {
            throw new Error('Missing required  parameter: tier');
        }

        if (parameters['tier'] !== undefined) {
            body = parameters['tier'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async deleteNetworksByNetworkIdTiersByTierIdImagesByImageName(
        parameters: {
            'networkId': string,
            'tierId': string,
            'imageName': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/images/{image_name}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['imageName'] === undefined) {
            throw new Error('Missing required  parameter: imageName');
        }

        path = path.replace('{image_name}', `${parameters['imageName']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getNetworksByNetworkIdTiersByTierIdName(
            parameters: {
                'networkId': string,
                'tierId': string,
            }
        ): Promise < tier_name >
        {
            let path = '/networks/{network_id}/tiers/{tier_id}/name';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['tierId'] === undefined) {
                throw new Error('Missing required  parameter: tierId');
            }

            path = path.replace('{tier_id}', `${parameters['tierId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdTiersByTierIdName(
        parameters: {
            'networkId': string,
            'tierId': string,
            'name': tier_name,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/name';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['name'] === undefined) {
            throw new Error('Missing required  parameter: name');
        }

        if (parameters['name'] !== undefined) {
            body = parameters['name'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdTiersByTierIdVersion(
            parameters: {
                'networkId': string,
                'tierId': string,
            }
        ): Promise < tier_version >
        {
            let path = '/networks/{network_id}/tiers/{tier_id}/version';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            if (parameters['tierId'] === undefined) {
                throw new Error('Missing required  parameter: tierId');
            }

            path = path.replace('{tier_id}', `${parameters['tierId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdTiersByTierIdVersion(
        parameters: {
            'networkId': string,
            'tierId': string,
            'version': tier_version,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/tiers/{tier_id}/version';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['tierId'] === undefined) {
            throw new Error('Missing required  parameter: tierId');
        }

        path = path.replace('{tier_id}', `${parameters['tierId']}`);

        if (parameters['version'] === undefined) {
            throw new Error('Missing required  parameter: version');
        }

        if (parameters['version'] !== undefined) {
            body = parameters['version'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getNetworksByNetworkIdType(
            parameters: {
                'networkId': string,
            }
        ): Promise < string >
        {
            let path = '/networks/{network_id}/type';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putNetworksByNetworkIdType(
        parameters: {
            'networkId': string,
            'type': string,
        }
    ): Promise < "Success" > {
        let path = '/networks/{network_id}/type';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['type'] === undefined) {
            throw new Error('Missing required  parameter: type');
        }

        if (parameters['type'] !== undefined) {
            body = parameters['type'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphony(): Promise < Array < string >
        >
        {
            let path = '/symphony';
            let body;
            let query = {};

            return await this.request(path, 'GET', query, body);
        }
    static async postSymphony(
        parameters: {
            'symphonyNetwork': symphony_network,
        }
    ): Promise < "Success" > {
        let path = '/symphony';
        let body;
        let query = {};
        if (parameters['symphonyNetwork'] === undefined) {
            throw new Error('Missing required  parameter: symphonyNetwork');
        }

        if (parameters['symphonyNetwork'] !== undefined) {
            body = parameters['symphonyNetwork'];
        }

        return await this.request(path, 'POST', query, body);
    }
    static async deleteSymphonyByNetworkId(
        parameters: {
            'networkId': string,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        return await this.request(path, 'DELETE', query, body);
    }
    static async getSymphonyByNetworkId(
            parameters: {
                'networkId': string,
            }
        ): Promise < symphony_network >
        {
            let path = '/symphony/{network_id}';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkId(
        parameters: {
            'networkId': string,
            'symphonyNetwork': symphony_network,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['symphonyNetwork'] === undefined) {
            throw new Error('Missing required  parameter: symphonyNetwork');
        }

        if (parameters['symphonyNetwork'] !== undefined) {
            body = parameters['symphonyNetwork'];
        }

        return await this.request(path, 'PUT', query, body);
    }
    static async getSymphonyByNetworkIdFeatures(
            parameters: {
                'networkId': string,
            }
        ): Promise < network_features >
        {
            let path = '/symphony/{network_id}/features';
            let body;
            let query = {};
            if (parameters['networkId'] === undefined) {
                throw new Error('Missing required  parameter: networkId');
            }

            path = path.replace('{network_id}', `${parameters['networkId']}`);

            return await this.request(path, 'GET', query, body);
        }
    static async putSymphonyByNetworkIdFeatures(
        parameters: {
            'networkId': string,
            'config': network_features,
        }
    ): Promise < "Success" > {
        let path = '/symphony/{network_id}/features';
        let body;
        let query = {};
        if (parameters['networkId'] === undefined) {
            throw new Error('Missing required  parameter: networkId');
        }

        path = path.replace('{network_id}', `${parameters['networkId']}`);

        if (parameters['config'] === undefined) {
            throw new Error('Missing required  parameter: config');
        }

        if (parameters['config'] !== undefined) {
            body = parameters['config'];
        }

        return await this.request(path, 'PUT', query, body);
    }
}
