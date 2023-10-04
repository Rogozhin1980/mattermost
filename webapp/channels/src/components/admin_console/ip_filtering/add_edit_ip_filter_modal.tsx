import React, { useState } from 'react';
import { Button, Modal } from 'react-bootstrap';
import Input, { CustomMessageInputType } from 'components/widgets/inputs/input/input';
import { AllowedIPRange } from '@mattermost/types/config';

import './add_edit_ip_filter_modal.scss'
import InfoIcon from 'components/widgets/icons/info_icon';
import { useIntl } from 'react-intl';

type Props = {
    onClose?: () => void;
    onSave?: (allowedIPRange: AllowedIPRange, oldIPRange?: AllowedIPRange) => void;
    existingRange?: AllowedIPRange;
    currentIP?: string;
}

function validateCIDR(cidr: string) {
    const cidrRegex = /^(\d{1,3}\.){3}\d{1,3}\/(1[0-9]|[1-9]|[0-2][0-9]|3[0-2])$/;
    return cidrRegex.test(cidr);
}

export default function IPFilteringAddOrEditModal({ onClose, onSave, existingRange, currentIP }: Props) {
    const {formatMessage} = useIntl();
    const [name, setName] = useState(existingRange?.Description || '');
    const [cidr, setCIDR] = useState(existingRange?.CIDRBlock || '');

    const [cidrError, setCidrError] = useState<CustomMessageInputType>(null);

    const handleSave = () => {
        const allowedIPRange: AllowedIPRange = {
            CIDRBlock: cidr,
            Description: name,
            Enabled: true,
            OwnerID: '',
        };

        if (existingRange) {
            onSave?.(allowedIPRange, existingRange);
        } else {
            onSave?.(allowedIPRange);
        }

        onClose?.();
    }

    const handleCIDRChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const cidr = e.target.value;
        setCIDR(cidr);
        setCidrError(null);
    }

    const validateCIDRInput = () => {
        if (!validateCIDR(cidr)) {
            setCidrError({ type: 'error', value: 'Invalid CIDR' });
        }
    }

    return (
        <Modal
            className={'IPFilteringAddOrEditModal'}
            dialogClassName={'IPFilteringAddOrEditModal__dialog'}
            show={true}
            onHide={() => onClose?.()}
        >
            <Modal.Header closeButton={true}>
                <div className='title'>
                    {formatMessage({id: 'admin.ip_filtering.add_ip_filter', defaultMessage: 'Add IP Filter'})}
                </div>
            </Modal.Header>
            <Modal.Body>
                <div className='body'>
                    <div className="current_ip_notice">
                        <div className="Content">
                            <span><InfoIcon />{formatMessage({id: 'admin.ip_filtering.your_current_ip_is', defaultMessage: 'Your current IP address is {ip}'}, {ip: currentIP})}</span>
                        </div>
                    </div>
                    <div className="inputs">
                        <div>
                            {formatMessage({ id: 'admin.ip_filtering.name', defaultMessage: 'Name' })}
                            <Input
                                type='text'
                                name='name'
                                onChange={(e) => setName(e.target.value)}
                                value={name}
                                placeholder={'Enter a name for this rule'}
                                required={true}
                                useLegend={false}
                            />
                        </div>
                        <div>{formatMessage({ id: 'admin.ip_filtering.allow_following_range', defaultMessage: 'Allow the following range of IP Addresses' })}
                        <Input
                            type='text'
                            name='ip_address_range'
                            onChange={handleCIDRChange}
                            onBlur={validateCIDRInput}
                            value={cidr}
                            placeholder={'Enter IP Range'}
                            required={true}
                            useLegend={false}
                            customMessage={cidrError}
                        />
                        </div>
                        {/* TODO: get proper PL for more info link out */}
                        <p>{formatMessage({id: 'admin.ip_filtering.more_info', defaultMessage: 'Enter ranges in CIDR format (e.g. 192.168.0.1/8). {link}'}, {link: <a href='https://docs.mattermost.com/deployment/ip-address-filtering.html' target='_blank' rel='noopener noreferrer'>{formatMessage({id: 'admin.ip_filtering.more_info_link', defaultMessage: 'More info'})}</a>})}</p>
                    </div>
                </div>
            </Modal.Body>
            <Modal.Footer>
                <Button
                    type="button"
                    className="btn-cancel"
                    onClick={() => onClose?.()}
                >
                    Cancel
                </Button>
                <Button
                    type="button"
                    className="btn-save"
                    onClick={handleSave}
                    disabled={cidrError !== null}
                >
                    {existingRange ? formatMessage({id: 'admin.ip_filtering.update_filter', defaultMessage: 'Update filter'}) : formatMessage({id: 'admin.ip_filtering.save', defaultMessage: 'Save'})}
                </Button>
            </Modal.Footer>
        </Modal>
    );
}