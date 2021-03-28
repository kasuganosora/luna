

function getPosts(){
    return [
        {
            id: 1,
            title: "网约车大数据杀熟收“苹果税”？复旦副教授花5万多元打车给答案",
            author: {id:123, name:"作者名字"},
            catalogs: [
                {id: "998", name: "分类1"},
                {id: "997", name: "分类2"}
            ],
            tags: [
                {id: "996", name: "标签1"},
                {id: "007", name: "标签2"},
            ],
            status: "draft",
            published_at: null,
            created_at: new Date()
        },
        {
            id: 2,
            title: "美网友围观SpaceX火箭爆炸现场 意外发现一只机器狗",
            author: {id:123, name:"作者名字"},
            catalogs: [
                {id: "998", name: "分类1"},
                {id: "997", name: "分类2"}
            ],
            tags: [
                {id: "996", name: "标签1"},
                {id: "007", name: "标签2"},
                {id: "8848", name: "黄金手机"},
            ],
            created_at: new Date(),
            status: "published",
            published_at: new Date("2017-01-02 12:31:59"),
        },
        {
            id: 3,
            title: "汉服竟如此赚钱？90后入坑花费数十万狂购700套，山东这个小镇赚翻了",
            author: {id:123, name:"作者名字"},
            catalogs: [
                {id: "998", name: "分类1"},
                {id: "997", name: "分类2"},
                {id: "233", name: "Q宝智能嘴炮"}
            ],
            tags: [
                {id: "996", name: "标签1"},
                {id: "007", name: "标签2"},
            ],
            created_at: new Date(),
            status: "draft",
            published_at: null,
        }
    ]
}

function getPost(){
    return {
        title: "标题是什么",
        markdown:"# 这是标题"
    }
}

export default {
    getPosts,
    getPost
}