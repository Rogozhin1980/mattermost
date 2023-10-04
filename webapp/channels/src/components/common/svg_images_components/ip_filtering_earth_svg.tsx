// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

type SvgProps = {
    width: number;
    height: number;
}

const IPFilteringEarthSvg = ({ width, height }: SvgProps) => (
    <svg width={width} height={height} viewBox="0 0 140 140" fill="none" xmlns="http://www.w3.org/2000/svg">
        <ellipse cx="69.5" cy="120" rx="29.5" ry="3" fill="black" fill-opacity="0.06" />
        <path d="M113.191 70.0004C113.19 78.1542 110.882 86.1411 106.533 93.0378C102.184 99.9345 95.9719 105.459 88.6152 108.974C81.2586 112.488 73.0577 113.848 64.9607 112.897C56.8636 111.945 49.2012 108.721 42.8592 103.598C36.5172 98.4739 31.7546 91.6597 29.1222 83.9427C26.4897 76.2257 26.0948 67.9212 27.9832 59.9892C29.8715 52.0572 33.9661 44.8217 39.7933 39.1192C45.6205 33.4168 52.9425 29.4802 60.9128 27.7647H62.1351L63.0116 27.3418C63.7189 27.2342 64.4339 27.1342 65.1489 27.0496L66.6865 27.0958L67.7397 26.8421H68.3548L70.0538 27.0804L71.1763 26.8267H71.8221L74.19 27.4264L76.9884 27.3495L77.9802 27.5187L88.2206 34.0467L99.3991 38.2987C103.764 42.3429 107.242 47.2479 109.616 52.7043C111.99 58.1607 113.208 64.0499 113.191 70.0004Z" fill="#1C58D9" />
        <path d="M72.0989 40.4363C74.9127 41.3359 74.7974 39.5598 76.5734 41.6281C77.5959 42.8199 79.441 44.9882 81.6321 45.5264C81.9396 45.6033 82.6238 44.9421 82.6238 44.4192C82.6238 43.8964 81.0863 42.7892 80.2406 42.5585C79.3949 42.3278 80.5865 40.6978 79.4256 40.198C78.2647 39.6982 76.9808 37.8451 76.5195 36.3919C76.0583 34.9387 74.0133 36.1612 73.3136 35.8691C72.614 35.5769 71.776 35.7537 71.0072 36.8609C70.2384 37.9682 71.3762 40.2057 72.0989 40.4363Z" fill="#FFBC1F" />
        <path d="M87.3369 81.1265C87.3369 80.5191 81.7708 77.4896 81.0942 77.4512C79.9871 77.382 78.4803 76.8284 77.4654 77.4512C77.0368 77.6945 76.6729 78.0374 76.4045 78.4508C76.2123 78.8122 75.1283 78.1432 74.7132 78.3662C73.7137 78.8967 73.1755 77.1667 73.4292 76.5055C74.1288 74.8831 71.661 75.9057 71.4073 75.4597C71.7071 75.9903 71.3227 70.8463 70.4616 73.7451C69.8082 75.9442 67.2481 73.4298 67.5556 71.8305C68.0322 69.3624 71.7532 69.001 72.5682 69.4393C72.9833 69.6622 74.4979 71.8536 74.6901 71.7921C75.1052 71.6614 74.4441 69.2086 74.3595 68.8856C74.0673 67.6862 76.3584 67.2863 76.5198 66.033C76.6582 64.9719 77.3809 65.341 78.0574 64.7797C79.0953 63.957 78.2804 62.4192 79.595 61.6349C80.0486 61.3658 81.7477 61.1736 81.9707 60.6968C82.1936 60.2201 80.5407 59.0745 81.9707 58.7361C83.2315 58.4286 84.5154 58.9745 84.8613 57.1522C85.1688 55.4991 83.3237 55.1223 82.7932 53.8382C82.6241 53.4538 82.4166 51.2394 81.6324 51.6469C80.6714 52.1543 80.0948 53.3769 79.2799 51.9237C77.8576 49.3017 74.3364 48.1714 75.9893 52.2927C76.835 54.4226 73.6829 60.9044 73.1986 54.9147C73.0371 52.831 66.8098 52.2082 68.901 49.9015C70.9921 47.5948 73.9213 47.4871 74.7055 43.8964C75.7972 38.9139 72.3452 43.1275 70.5846 42.1894C67.9553 40.7516 70.1849 41.5281 68.8241 43.3889C68.309 44.0963 66.6176 42.9429 65.8642 43.0583C65.1108 43.1736 63.3809 44.7037 62.8966 43.7426C62.0202 41.9818 55.8467 40.8976 54.1784 41.5205C50.6957 42.8276 49.4887 41.0591 46.0137 40.3748C44.8067 40.1364 38.8792 39.9673 39.2636 42.0202C39.5404 42.7727 39.864 43.5071 40.2323 44.2193C40.609 45.3496 38.9099 44.4807 38.5179 44.7498C39.21 45.3637 39.9337 45.9412 40.6859 46.4799C41.2548 47.0642 38.5102 47.7716 38.2334 48.7404C37.6875 50.7088 42.0082 50.1475 41.2164 52.5157C40.9703 53.2308 38.5025 54.4072 38.2641 56.8293C38.218 57.2906 41.3393 53.8306 43.3536 54.5226C45.1295 55.1685 44.9758 50.3398 47.1976 50.3398C51.3184 50.3398 52.118 54.3534 54.286 57.0446C55.3008 58.2979 56.1158 58.8207 56.3003 60.4585C56.4848 62.0963 56.1388 63.8416 56.6924 65.2949C57.2459 66.7481 58.2915 66.7789 58.9988 67.9014C59.8445 69.2163 60.3211 70.5003 61.4666 71.646C63.1272 73.3145 65.9564 75.8673 68.5704 75.4905C69.0855 75.4136 71.0306 76.5208 71.4688 76.8514C72.0759 77.4221 72.5941 78.0804 73.0064 78.8045C73.4677 79.4196 74.9284 79.0351 75.5665 79.3273C76.2046 79.6195 74.0827 83.9407 74.8515 85.0479C75.6203 86.1551 75.6895 87.7006 76.9427 88.6848C78.1958 89.669 78.4803 89.6921 78.6417 91.3145C78.8801 93.9441 78.1112 96.6045 77.4808 99.0573C77.1502 100.311 77.8114 101.702 77.3194 102.902C76.7043 104.401 76.2892 105.693 76.8274 107.361C77.7576 110.26 82.9778 110.729 79.9026 108.2C78.5879 107.123 81.3249 104.524 80.4715 104.278C79.0031 103.663 80.7405 103.109 80.9328 102.471C81.3248 101.203 87.1447 95.951 87.5752 95.0821C87.9366 94.3516 87.7059 93.5981 88.2979 92.9676C88.8899 92.3371 89.8355 92.514 90.6043 92.0219C92.1419 90.9915 91.5884 89.2461 92.2265 87.8237C93.4873 85.0018 94.0178 86.1167 92.288 83.91C91.6345 83.1103 87.4522 83.4025 87.3369 81.1265Z" fill="#FFBC1F" />
        <path d="M56.2 38.4295C57.261 39.5367 58.3142 37.8528 58.5064 38.4295C58.6986 39.0062 58.714 40.1134 59.9748 40.198C61.2357 40.2826 58.7986 41.1053 59.6135 41.4513C60.4284 41.7973 62.1736 42.9045 62.8117 42.4432C63.4498 41.9818 64.2109 41.6282 64.4954 41.9742C64.7798 42.3202 66.5327 41.9203 66.033 41.336C65.5333 40.7516 64.2032 39.3753 64.2571 38.5833C64.3109 37.7913 65.7947 35.7537 64.1418 35.9306C63.743 35.9807 63.3796 36.185 63.1296 36.4998C62.8796 36.8145 62.7627 37.2146 62.804 37.6145C62.804 38.1989 61.2664 36.7456 60.9435 36.7687C60.6206 36.7918 59.5443 35.0002 59.1368 35.0002C58.7294 35.0002 56.0001 34.1237 56.2307 35.0002C56.4614 35.8768 55.1391 37.3223 56.2 38.4295Z" fill="#FFBC1F" />
        <path d="M80.6022 46.749C79.7796 46.28 78.2189 44.1347 77.527 44.8268C76.8351 45.5188 77.7038 47.1335 78.2958 47.7947C78.8878 48.456 81.3248 50.0092 80.9712 49.3325C80.6175 48.6559 82.2167 47.6794 80.6022 46.749Z" fill="#FFBC1F" />
        <path d="M57.8068 31.9554C58.2297 31.9169 58.8909 33.3855 59.5136 33.547C60.1363 33.7085 60.4439 34.4389 60.6745 34.6081C60.9051 34.7772 63.5114 33.7777 64.1264 33.547C64.7415 33.3163 63.3576 31.8785 62.8425 30.8712C62.3274 29.864 62.7272 32.5397 61.9968 32.1015C61.2665 31.6632 61.3741 31.0634 60.4592 31.0634C59.5443 31.0634 60.1824 29.2412 59.06 29.6256C57.9375 30.0101 57.384 31.9938 57.8068 31.9554Z" fill="#FFBC1F" />
        <path d="M67.1554 35.4462C66.3866 35.5999 66.5326 36.0997 66.1021 36.5303C65.9011 36.7358 65.7886 37.0118 65.7886 37.2992C65.7886 37.5866 65.9011 37.8626 66.1021 38.0681C66.4942 38.5679 67.1092 39.4598 67.5782 38.6063C68.0472 37.7529 67.5782 37.0685 67.5782 36.5149C67.5782 35.9613 67.5321 35.3693 67.1554 35.4462Z" fill="#FFBC1F" />
        <path d="M68.7011 38.1604C68.8933 38.6218 69.2469 37.3454 70.0003 36.6841C70.7538 36.0229 69.6313 34.3774 68.9317 35.2155C68.4627 35.7922 68.5089 37.6914 68.7011 38.1604Z" fill="#FFBC1F" />
        <path d="M68.5626 38.737C68.1628 39.306 67.8399 40.8284 68.5626 40.9514C68.7209 40.9697 68.8811 40.9437 69.0256 40.8763C69.17 40.8089 69.2928 40.7028 69.3805 40.5697C69.4681 40.4366 69.5172 40.2818 69.5221 40.1225C69.527 39.9632 69.4876 39.8056 69.4083 39.6674C69.07 39.0907 68.7317 38.5063 68.5626 38.737Z" fill="#FFBC1F" />
        <path d="M71.4916 28.9256C70.7996 29.2409 72.5064 29.8176 71.3378 30.1328C70.1692 30.4481 70.1385 31.4092 69.1313 30.5941C68.1242 29.7791 67.0248 30.8709 67.7859 31.363C68.547 31.8551 69.8232 33.07 69.7848 33.7313C69.7464 34.3925 70.9918 34.7001 72.1527 34.6232C73.3136 34.5463 74.7129 34.8923 74.7513 34.3925C74.7897 33.8927 74.9281 33.6775 74.7897 33.3007C74.5898 32.7548 74.0901 32.2165 72.8908 32.6856C71.6914 33.1546 70.9072 32.9162 70.7381 32.3319C70.569 31.7475 73.0907 30.6249 74.2054 31.4015C75.3202 32.1781 76.3812 30.9786 75.5662 30.4327C74.7513 29.8868 75.2971 29.3101 75.8584 29.2332C76.2658 29.1794 76.6271 28.149 76.9347 27.3494C75.2246 27.0728 73.4993 26.9007 71.7683 26.8342C71.776 27.9338 72.1835 28.6104 71.4916 28.9256Z" fill="#FFBC1F" />
        <path d="M61.1819 27.957C61.9507 28.2877 63.227 27.957 62.9886 27.365C62.289 27.488 61.5894 27.6187 60.8975 27.7648C60.9744 27.852 61.0723 27.9182 61.1819 27.957Z" fill="#FFBC1F" />
        <path d="M61.6585 28.6107C60.5668 28.0494 60.6283 29.4257 61.128 29.9716C61.6277 30.5175 63.0731 29.3027 61.6585 28.6107Z" fill="#FFBC1F" />
        <path d="M68.3013 26.8345C68.6367 26.895 68.9503 27.0427 69.2108 27.2627C69.4712 27.4826 69.6692 27.7672 69.7851 28.0878C69.9311 28.5799 71.1074 27.5111 71.1228 26.8114H70.0003C69.4314 26.8037 68.8702 26.8114 68.3013 26.8345Z" fill="#FFBC1F" />
        <path d="M65.9942 28.4492C65.5022 29.1258 65.925 29.9408 66.5093 29.2642C67.0936 28.5876 66.6246 26.8729 67.0859 27.6111C67.5472 28.3492 68.4083 29.5179 68.9849 29.5718C69.5615 29.6256 69.0618 28.8874 68.9849 28.2723C68.908 27.6572 67.5319 27.0728 67.6856 26.8652C66.8246 26.9114 65.9712 26.9806 65.1255 27.0728C65.7328 27.4727 66.371 27.9648 65.9942 28.4492Z" fill="#FFBC1F" />
        <path d="M68.5625 34.2775C69.2467 34.5543 68.5625 32.8858 68.5625 32.8858C67.8475 33.3087 67.8859 33.993 68.5625 34.2775Z" fill="#FFBC1F" />
        <path d="M66.1024 33.224C66.2177 33.693 66.9558 33.8468 67.3478 33.224C67.7399 32.6012 67.3478 30.3406 66.7943 31.0557C66.2407 31.7708 65.6026 31.4401 64.4879 31.0557C63.3731 30.6712 65.9948 32.7319 66.1024 33.224Z" fill="#FFBC1F" />
        <path d="M73.6826 72.8378C74.6359 72.7225 75.2202 73.8835 76.2966 73.8297C77.3729 73.7759 79.0027 75.4059 79.933 74.4371C80.7018 73.6682 76.6118 73.0147 75.3202 72.315C74.0286 71.6153 72.737 72.9531 73.6826 72.8378Z" fill="#FFBC1F" />
        <path d="M78.4337 29.1181C79.7714 29.7639 79.5408 29.7409 78.7259 30.3637C77.9109 30.9865 78.2031 32.3398 80.0021 32.2475C81.8011 32.1552 82.9159 30.3637 83.7308 32.8857C84.5457 35.4077 83.9614 38.1603 85.007 38.9369C86.0526 39.7135 85.1838 40.2364 85.007 40.9591C84.8302 41.6819 85.2991 42.2047 86.1141 41.9741C86.929 41.7434 86.6984 43.2581 86.1141 43.8963C85.5298 44.5345 84.8379 45.7571 86.6369 47.9715C88.4359 50.1859 88.6742 49.8861 89.2508 50.6473C89.8274 51.4085 90.3041 50.3628 90.4732 49.4862C90.6424 48.6097 90.0657 45.5264 92.1338 45.5264C92.778 45.5195 93.4029 45.3057 93.9164 44.9167C94.4299 44.5277 94.8049 43.984 94.9861 43.3658C95.2167 42.8429 96.8466 42.8968 97.7153 42.1432C98.5841 41.3897 99.3529 40.9053 99.3759 38.2987C93.4007 32.7528 85.9725 29.0196 77.9571 27.5341C77.7649 28.1339 77.7726 28.8105 78.4337 29.1181Z" fill="#FFBC1F" />
    </svg>

)

export default IPFilteringEarthSvg;