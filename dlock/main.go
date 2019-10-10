package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type EtcdMutex struct {
	Ttl     int64              //租约时间
	Conf    clientv3.Config    //etcd集群配置
	Key     string             //etcd的key
	cancel  context.CancelFunc //关闭续租的func
	lease   clientv3.Lease
	leaseID clientv3.LeaseID
	txn     clientv3.Txn
}

func (em *EtcdMutex) init() error {
	var err error
	var ctx context.Context
	client, err := clientv3.New(em.Conf)
	if err != nil {
		return err
	}
	em.txn = clientv3.NewKV(client).Txn(context.TODO())
	em.lease = clientv3.NewLease(client)
	leaseResp, err := em.lease.Grant(context.TODO(), em.Ttl)
	if err != nil {
		return err
	}
	ctx, em.cancel = context.WithCancel(context.TODO())
	em.leaseID = leaseResp.ID
	_, err = em.lease.KeepAlive(ctx, em.leaseID)
	return err
}

func (em *EtcdMutex) Lock() error {
	err := em.init()
	if err != nil {
		return err
	}
	//LOCK:
	em.txn.If(clientv3.Compare(clientv3.CreateRevision(em.Key), "=", 0)).
		Then(clientv3.OpPut(em.Key, "", clientv3.WithLease(em.leaseID))).
		Else()
	txnResp, err := em.txn.Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded { //判断txn.if条件是否成立
		return fmt.Errorf("抢锁失败")
	}
	return nil
}

func (em *EtcdMutex) UnLock() {
	em.cancel()
	em.lease.Revoke(context.TODO(), em.leaseID)
	fmt.Println("释放了锁")
}

func main() {
	var conf = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379", "127.0.0.1:2380"},
		DialTimeout: 5 * time.Second,
	}
	eMutex1 := &EtcdMutex{
		Conf: conf,
		Ttl:  10,
		Key:  "lock",
	}
	eMutex2 := &EtcdMutex{
		Conf: conf,
		Ttl:  10,
		Key:  "lock",
	}
	//groutine1
	go func() {
		err := eMutex1.Lock()
		if err != nil {
			fmt.Println("groutine1抢锁失败")
			fmt.Println(err)
			return
		}
		fmt.Println("groutine1抢锁成功")
		time.Sleep(10 * time.Second)
		defer eMutex1.UnLock()
	}()

	//groutine2
	go func() {
		err := eMutex2.Lock()
		if err != nil {
			fmt.Println("groutine2抢锁失败")
			fmt.Println(err)
			return
		}
		fmt.Println("groutine2抢锁成功")
		defer eMutex2.UnLock()
	}()
	time.Sleep(30 * time.Second)
}
